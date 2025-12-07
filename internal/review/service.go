package review

import (
	"strings"
	"time"

	"github.com/Gab-Mello/service-finder/internal/order"
)

type Service struct {
	repo    Repository
	orders  order.Repository
	now     func() time.Time
	editTTL time.Duration
}

func NewService(r Repository, orders order.Repository, now func() time.Time) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{
		repo:    r,
		orders:  orders,
		now:     now,
		editTTL: 24 * time.Hour,
	}
}

func (s *Service) Create(clientID, orderID string, stars int, comment string) (*Review, error) {
	if stars < 1 || stars > 5 {
		return nil, ErrInvalidFields
	}
	ord, err := s.orders.ByID(orderID)
	if err != nil {
		return nil, err
	}
	if ord.Status != order.StatusCompleted {
		return nil, ErrOrderNotDone
	}
	if ord.ClientID != clientID {
		return nil, ErrForbidden
	}
	if _, err := s.repo.ByOrderID(orderID); err == nil {
		return nil, ErrAlreadyExists
	}

	now := s.now()
	rv := &Review{
		OrderID:    orderID,
		ClientID:   clientID,
		ProviderID: ord.ProviderID,
		Stars:      stars,
		Comment:    strings.TrimSpace(comment),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.repo.Create(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

func (s *Service) Edit(clientID, orderID string, stars int, comment string) (*Review, error) {
	rv, err := s.repo.ByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	if rv.ClientID != clientID {
		return nil, ErrForbidden
	}
	if s.now().After(rv.CreatedAt.Add(s.editTTL)) {
		return nil, ErrEditWindowOver
	}
	if stars < 1 || stars > 5 {
		return nil, ErrInvalidFields
	}
	rv.Stars = stars
	rv.Comment = strings.TrimSpace(comment)
	rv.UpdatedAt = s.now()
	if err := s.repo.Update(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

func (s *Service) AvgForProvider(providerID string) (avg float64, count int) {
	list, _ := s.repo.ListByProvider(providerID)
	if len(list) == 0 {
		return 0, 0
	}
	var sum int
	for _, r := range list {
		sum += r.Stars
	}
	return float64(sum) / float64(len(list)), len(list)
}
