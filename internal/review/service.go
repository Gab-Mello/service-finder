package review

import (
	"log"
	"strings"
	"time"

	"github.com/Gab-Mello/service-finder/internal/order"
)

const maxCommentLen = 2000

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
	comment = strings.TrimSpace(comment)
	if len(comment) > maxCommentLen {
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

	_, err = s.repo.ByOrderID(orderID)
	if err == nil {
		return nil, ErrAlreadyExists
	}
	if err != ErrNotFound {
		log.Printf("error checking existing review for order %s: %v", orderID, err)
		return nil, err
	}

	now := s.now()
	rv := &Review{
		OrderID:    orderID,
		ClientID:   clientID,
		ProviderID: ord.ProviderID,
		Stars:      stars,
		Comment:    comment,
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
	comment = strings.TrimSpace(comment)
	if len(comment) > maxCommentLen {
		return nil, ErrInvalidFields
	}

	rv.Stars = stars
	rv.Comment = comment
	rv.UpdatedAt = s.now()
	if err := s.repo.Update(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

func (s *Service) ByOrderID(orderID string) (*Review, error) {
	return s.repo.ByOrderID(orderID)
}

func (s *Service) AvgForProvider(providerID string) (avg float64, count int) {
	list, err := s.repo.ListByProvider(providerID)
	if err != nil {
		log.Printf("error listing reviews for provider %s: %v", providerID, err)
		return 0, 0
	}
	if len(list) == 0 {
		return 0, 0
	}
	var sum int
	for _, r := range list {
		sum += r.Stars
	}
	return float64(sum) / float64(len(list)), len(list)
}
