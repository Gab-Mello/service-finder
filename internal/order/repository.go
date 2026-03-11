package order

import "sync"

type Repository interface {
	Create(o *Order) error
	ByID(id string) (*Order, error)
	Update(o *Order) error
	ListMine(userID string) ([]Order, error)
}

type memoryRepo struct {
	mu     sync.RWMutex
	byID   map[string]*Order
	byUser map[string][]string // userID -> []orderID index (for both client and provider)
}

func NewRepository() Repository {
	return &memoryRepo{
		byID:   make(map[string]*Order),
		byUser: make(map[string][]string),
	}
}

func (r *memoryRepo) Create(o *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c := *o
	r.byID[o.ID] = &c

	r.byUser[o.ClientID] = append(r.byUser[o.ClientID], o.ID)
	if o.ProviderID != o.ClientID {
		r.byUser[o.ProviderID] = append(r.byUser[o.ProviderID], o.ID)
	}
	return nil
}

func (r *memoryRepo) ByID(id string) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	c := *o
	return &c, nil
}

func (r *memoryRepo) Update(o *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[o.ID]; !ok {
		return ErrNotFound
	}
	c := *o
	r.byID[o.ID] = &c
	return nil
}

func (r *memoryRepo) ListMine(userID string) ([]Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[userID]
	out := make([]Order, 0, len(ids))
	for _, id := range ids {
		if o, ok := r.byID[id]; ok {
			out = append(out, *o)
		}
	}
	return out, nil
}
