package review

// memória simples (sem locks) para MVP
type Repository interface {
	Create(r *Review) error
	ByOrderID(orderID string) (*Review, error)
	Update(r *Review) error
	ListByProvider(providerID string) ([]Review, error)
}

type memoryRepo struct {
	byOrder    map[string]*Review
	byProvider map[string][]*Review
}

func NewRepository() Repository {
	return &memoryRepo{
		byOrder:    map[string]*Review{},
		byProvider: map[string][]*Review{},
	}
}

func (r *memoryRepo) Create(rv *Review) error {
	if _, ok := r.byOrder[rv.OrderID]; ok {
		return ErrAlreadyExists
	}
	c := *rv
	r.byOrder[rv.OrderID] = &c
	r.byProvider[rv.ProviderID] = append(r.byProvider[rv.ProviderID], &c)
	return nil
}

func (r *memoryRepo) ByOrderID(orderID string) (*Review, error) {
	v, ok := r.byOrder[orderID]
	if !ok {
		return nil, ErrNotFound
	}
	c := *v
	return &c, nil
}

func (r *memoryRepo) Update(rv *Review) error {
	old, ok := r.byOrder[rv.OrderID]
	if !ok {
		return ErrNotFound
	}
	// atualizar entry de byOrder
	*old = *rv

	// atualizar array do provider (substituir a que tem o mesmo OrderID)
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
	arr := r.byProvider[providerID]
	out := make([]Review, 0, len(arr))
	for _, p := range arr {
		out = append(out, *p)
	}
	return out, nil
}
