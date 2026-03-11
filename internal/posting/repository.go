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
	mu         sync.RWMutex
	byID       map[string]Posting
	byProvider map[string][]string // providerID -> []postingID index
}

func NewRepository() Repository {
	return &memoryRepo{
		byID:       make(map[string]Posting),
		byProvider: make(map[string][]string),
	}
}

func (r *memoryRepo) Create(p *Posting) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byID[p.ID]; exists {
		return ErrInvalidFields // ID already exists
	}

	r.byID[p.ID] = *p
	r.byProvider[p.ProviderID] = append(r.byProvider[p.ProviderID], p.ID)
	return nil
}

func (r *memoryRepo) Update(p *Posting) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[p.ID]; !ok {
		return ErrNotFound
	}
	r.byID[p.ID] = *p
	return nil
}

func (r *memoryRepo) ByID(id string) (*Posting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	it, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &it, nil
}

func (r *memoryRepo) ListByProvider(pid string) ([]Posting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byProvider[pid]
	out := make([]Posting, 0, len(ids))
	for _, id := range ids {
		if p, ok := r.byID[id]; ok {
			out = append(out, p)
		}
	}
	return out, nil
}

func (r *memoryRepo) ListPublic() ([]Posting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Posting, 0)
	for _, it := range r.byID {
		if !it.Archived {
			out = append(out, it)
		}
	}
	return out, nil
}
