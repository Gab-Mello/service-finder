package posting

import (
	"errors"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Gab-Mello/service-finder/internal/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo      Repository
	providers ports.ProviderDirectory
	ratings   ports.Ratings
	now       func() time.Time
	idgen     func() string
}

func NewService(r Repository, providers ports.ProviderDirectory, now func() time.Time, idgen func() string, ratings ports.Ratings) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	if idgen == nil {
		idgen = func() string { return uuid.NewString() }
	}
	return &Service{
		repo:      r,
		providers: providers,
		ratings:   ratings,
		now:       now,
		idgen:     idgen,
	}
}

func (s *Service) Create(providerID, title, desc string, price int64, category, city, district string) (*Posting, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(desc) == "" || price <= 0 ||
		strings.TrimSpace(category) == "" || strings.TrimSpace(city) == "" || strings.TrimSpace(district) == "" {
		return nil, ErrInvalidFields
	}

	providerName, err := s.providers.GetNameByID(providerID)
	if err != nil {
		return nil, ErrInvalidFields
	}

	p := &Posting{
		ID: s.idgen(), ProviderID: providerID, ProviderName: providerName,
		Title: strings.TrimSpace(title), Description: strings.TrimSpace(desc),
		Price:    price,
		Category: strings.TrimSpace(category),
		City:     strings.TrimSpace(city), District: strings.TrimSpace(district),
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
	s.enrich(p)
	return p, nil
}

func (s *Service) ListMine(providerID string) ([]Posting, error) {
	list, err := s.repo.ListByProvider(providerID)
	if err == nil {
		s.enrichMany(list)
	}
	return list, err
}

func (s *Service) ListPublic() ([]Posting, error) {
	list, err := s.repo.ListPublic()
	if err == nil {
		s.enrichMany(list)
	}
	return list, err
}

type SearchParams struct {
	Query                    string
	Category, City, District string
	PriceMin, PriceMax       int64

	RatingMin float64
	Sort      string
	Order     string
	Limit     int
	Offset    int
}

func (s *Service) Search(p SearchParams) ([]Posting, int) {
	all, _ := s.repo.ListPublic()

	if u, err := url.QueryUnescape(p.Query); err == nil {
		p.Query = u
	}
	if u, err := url.QueryUnescape(p.Category); err == nil {
		p.Category = u
	}
	if u, err := url.QueryUnescape(p.City); err == nil {
		p.City = u
	}
	if u, err := url.QueryUnescape(p.District); err == nil {
		p.District = u
	}

	norm := func(x string) string {
		x = strings.TrimSpace(strings.ToLower(x))
		return strings.Join(strings.Fields(x), " ")
	}

	q := norm(p.Query)
	wantCat := norm(p.Category)
	wantCity := norm(p.City)
	wantDist := norm(p.District)

	filtered := make([]Posting, 0, len(all))
	for _, it := range all {
		titleDesc := norm(it.Title + " " + it.Description)
		gotCat := norm(it.Category)
		gotCity := norm(it.City)
		gotDist := norm(it.District)

		if q != "" && !strings.Contains(titleDesc, q) {
			continue
		}
		if wantCat != "" && gotCat != wantCat {
			continue
		}
		if wantCity != "" && gotCity != wantCity {
			continue
		}
		if wantDist != "" && gotDist != wantDist {
			continue
		}
		if p.PriceMin > 0 && it.Price < p.PriceMin {
			continue
		}
		if p.PriceMax > 0 && it.Price > p.PriceMax {
			continue
		}

		filtered = append(filtered, it)
	}

	sortKey := strings.ToLower(p.Sort)
	order := strings.ToLower(p.Order)
	if sortKey == "" {
		sortKey = "relevance"
	}
	less := func(i, j int) bool {
		switch sortKey {
		case "price":
			if order == "desc" {
				return filtered[i].Price > filtered[j].Price
			}
			return filtered[i].Price < filtered[j].Price
		case "rating":
			fallthrough
		default:
			qi := 0
			qj := 0
			if q != "" && strings.Contains(norm(filtered[i].Title), q) {
				qi = 1
			}
			if q != "" && strings.Contains(norm(filtered[j].Title), q) {
				qj = 1
			}
			if qi != qj {
				return qi > qj
			}

			return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
		}
	}
	sort.SliceStable(filtered, less)

	limit := p.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := p.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= len(filtered) {
		return []Posting{}, -1
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	next := -1
	if end < len(filtered) {
		next = end
	}

	page := filtered[offset:end]
	s.enrichMany(page)

	return page, next
}

func (s *Service) enrich(p *Posting) {
	if s.ratings == nil || p == nil {
		return
	}
	if avg, _ := s.ratings.AvgForProvider(p.ProviderID); avg > 0 {
		p.ProviderAvg = avg
	}
}
func (s *Service) enrichMany(list []Posting) {
	if s.ratings == nil {
		return
	}
	for i := range list {
		if avg, _ := s.ratings.AvgForProvider(list[i].ProviderID); avg > 0 {
			list[i].ProviderAvg = avg
		}
	}
}
