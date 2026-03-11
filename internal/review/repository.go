package review

import "sync"

type Repository interface {
	Create(r *Review) error
	ByOrderID(orderID string) (*Review, error)
	Update(r *Review) error
	ListByProvider(providerID string) ([]Review, error)
}

type memoryRepo struct {
	mu         sync.RWMutex
	byOrder    map[string]*Review
	byProvider map[string][]*Review
}

func NewRepository() Repository {
	return &memoryRepo{
		byOrder:    make(map[string]*Review),
		byProvider: make(map[string][]*Review),
	}
}

func (r *memoryRepo) Create(rv *Review) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byOrder[rv.OrderID]; ok {
		return ErrAlreadyExists
	}
	c := *rv
	r.byOrder[rv.OrderID] = &c
	r.byProvider[rv.ProviderID] = append(r.byProvider[rv.ProviderID], &c)
	return nil
}

func (r *memoryRepo) ByOrderID(orderID string) (*Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.byOrder[orderID]
	if !ok {
		return nil, ErrNotFound
	}
	c := *v
	return &c, nil
}

func (r *memoryRepo) Update(rv *Review) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	old, ok := r.byOrder[rv.OrderID]
	if !ok {
		return ErrNotFound
	}

	*old = *rv

	arr := r.byProvider[rv.ProviderID]
	for i := range arr {
		if arr[i].OrderID == rv.OrderID {
			*(arr[i]) = *rv
			break
		}
	}
	return nil
}

func (r *memoryRepo) ListByProvider(providerID string) ([]Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	arr := r.byProvider[providerID]
	out := make([]Review, 0, len(arr))
	for _, p := range arr {
		out = append(out, *p)
	}
	return out, nil
}
