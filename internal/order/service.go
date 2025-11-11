package order

import (
	"time"

	"github.com/google/uuid"
)

type Notifier interface {
	OrderStatusChanged(o *Order)
}

type noopNotifier struct{}

func (noopNotifier) OrderStatusChanged(o *Order) {}

type Service struct {
	repo     Repository
	now      func() time.Time
	idgen    func() string
	notifier Notifier
}

func NewService(r Repository, now func() time.Time, idgen func() string, n Notifier) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	if idgen == nil {
		idgen = func() string { return uuid.NewString() }
	}
	if n == nil {
		n = noopNotifier{}
	}
	return &Service{repo: r, now: now, idgen: idgen, notifier: n}
}

func (s *Service) Request(clientID, postingID, providerID string) (*Order, error) {
	if clientID == "" || postingID == "" || providerID == "" {
		return nil, ErrInvalidFields
	}
	now := s.now()
	o := &Order{
		ID:         s.idgen(),
		PostingID:  postingID,
		ClientID:   clientID,
		ProviderID: providerID,
		Status:     StatusPending,
		History: []HistoryEntry{{
			At: now, By: clientID, From: "", To: StatusPending, Note: "pedido criado",
		}},
		CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.Create(o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) Accept(providerID, orderID string, scheduled time.Time) (*Order, error) {
	o, err := s.repo.ByID(orderID)
	if err != nil {
		return nil, err
	}
	if o.ProviderID != providerID {
		return nil, ErrForbidden
	}
	if o.Status != StatusPending {
		return nil, ErrInvalidState
	}
	o.ScheduledAt = &scheduled
	s.transition(o, providerID, StatusAccepted, "aceito com data/horário")
	return o, s.repo.Update(o)
}

func (s *Service) Start(providerID, orderID string) (*Order, error) {
	o, err := s.repo.ByID(orderID)
	if err != nil {
		return nil, err
	}
	if o.ProviderID != providerID {
		return nil, ErrForbidden
	}
	if o.Status != StatusAccepted {
		return nil, ErrInvalidState
	}
	s.transition(o, providerID, StatusInProgress, "início do serviço")
	return o, s.repo.Update(o)
}

func (s *Service) Complete(providerID, orderID string) (*Order, error) {
	o, err := s.repo.ByID(orderID)
	if err != nil {
		return nil, err
	}
	if o.ProviderID != providerID {
		return nil, ErrForbidden
	}
	if o.Status != StatusInProgress {
		return nil, ErrInvalidState
	}
	s.transition(o, providerID, StatusCompleted, "serviço concluído")
	return o, s.repo.Update(o)
}

func (s *Service) Cancel(actorID, orderID string) (*Order, error) {
	o, err := s.repo.ByID(orderID)
	if err != nil {
		return nil, err
	}
	if actorID != o.ClientID && actorID != o.ProviderID {
		return nil, ErrForbidden
	}
	if !(o.Status == StatusPending || o.Status == StatusAccepted) {
		return nil, ErrInvalidState
	}
	s.transition(o, actorID, StatusCanceled, "cancelado")
	return o, s.repo.Update(o)
}

func (s *Service) Get(id string) (*Order, error) { return s.repo.ByID(id) }

func (s *Service) ListMine(userID string) ([]Order, error) {
	return s.repo.ListMine(userID)
}

func (s *Service) transition(o *Order, by string, to Status, note string) {
	now := s.now()
	from := o.Status
	o.Status = to
	o.UpdatedAt = now
	o.History = append(o.History, HistoryEntry{At: now, By: by, From: from, To: to, Note: note})
	s.notifier.OrderStatusChanged(o)
}
