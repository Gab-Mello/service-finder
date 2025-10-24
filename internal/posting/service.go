package posting

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo  Repository
	now   func() time.Time
	idgen func() string
}

func NewService(r Repository, now func() time.Time, idgen func() string) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	if idgen == nil {
		idgen = func() string { return uuid.NewString() }
	}
	return &Service{repo: r, now: now, idgen: idgen}
}

func (s *Service) Create(providerID, title, desc string, price int64, category, city, district string) (*Posting, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(desc) == "" || price <= 0 ||
		strings.TrimSpace(category) == "" || strings.TrimSpace(city) == "" || strings.TrimSpace(district) == "" {
		return nil, ErrInvalidFields
	}
	p := &Posting{
		ID: s.idgen(), ProviderID: providerID,
		Title: strings.TrimSpace(title), Description: strings.TrimSpace(desc),
		Price: price, Category: strings.TrimSpace(category),
		City: strings.TrimSpace(city), District: strings.TrimSpace(district),
		CreatedAt: s.now(), UpdatedAt: s.now(),
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Update(providerID, id string, patch map[string]any) (*Posting, error) {
	p, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	if p.ProviderID != providerID {
		return nil, ErrForbidden
	}

	if v, ok := patch["title"].(string); ok {
		p.Title = strings.TrimSpace(v)
	}
	if v, ok := patch["description"].(string); ok {
		p.Description = strings.TrimSpace(v)
	}
	if v, ok := patch["category"].(string); ok {
		p.Category = strings.TrimSpace(v)
	}
	if v, ok := patch["city"].(string); ok {
		p.City = strings.TrimSpace(v)
	}
	if v, ok := patch["district"].(string); ok {
		p.District = strings.TrimSpace(v)
	}
	if v, ok := patch["price"].(float64); ok {
		if v <= 0 {
			return nil, errors.New("price must be > 0")
		}
		p.Price = int64(v)
	}

	p.UpdatedAt = s.now()
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Archive(providerID, id string) error {
	p, err := s.repo.ByID(id)
	if err != nil {
		return err
	}
	if p.ProviderID != providerID {
		return ErrForbidden
	}
	p.Archived = true
	p.UpdatedAt = s.now()
	return s.repo.Update(p)
}

func (s *Service) GetPublic(id string) (*Posting, error) {
	p, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	if p.Archived {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *Service) ListMine(providerID string) ([]Posting, error) {
	return s.repo.ListByProvider(providerID)
}
func (s *Service) ListPublic() ([]Posting, error) { return s.repo.ListPublic() }
