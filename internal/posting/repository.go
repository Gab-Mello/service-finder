package posting

import "sync"

type Repository interface {
	Create(*Posting) error
	Update(*Posting) error
	ByID(id string) (*Posting, error)
	ListByProvider(providerID string) ([]Posting, error)
	ListPublic() ([]Posting, error)
}

type memoryRepo struct {
	mu   sync.RWMutex
	byID map[string]*Posting
}

func NewRepository() Repository {
	return &memoryRepo{byID: map[string]*Posting{}}
}

func (r *memoryRepo) Create(p *Posting) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	r.byID[p.ID] = &cp
	return nil
}

func (r *memoryRepo) Update(p *Posting) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[p.ID]; !ok {
		return ErrNotFound
	}
	cp := *p
	r.byID[p.ID] = &cp
	return nil
}

func (r *memoryRepo) ByID(id string) (*Posting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *memoryRepo) ListByProvider(pid string) ([]Posting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Posting{}
	for _, p := range r.byID {
		if p.ProviderID == pid {
			out = append(out, *p)
		}
	}
	return out, nil
}

func (r *memoryRepo) ListPublic() ([]Posting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Posting{}
	for _, p := range r.byID {
		if !p.Archived {
			out = append(out, *p)
		}
	}
	return out, nil
}
