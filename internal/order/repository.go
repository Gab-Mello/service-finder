package order

type Repository interface {
	Create(o *Order) error
	ByID(id string) (*Order, error)
	Update(o *Order) error
	ListMine(userID string) ([]Order, error)
}

type memoryRepo struct {
	byID map[string]*Order
}

func NewRepository() Repository {
	return &memoryRepo{byID: map[string]*Order{}}
}

func (r *memoryRepo) Create(o *Order) error {
	c := *o
	r.byID[o.ID] = &c
	return nil
}

func (r *memoryRepo) ByID(id string) (*Order, error) {
	o, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	c := *o
	return &c, nil
}

func (r *memoryRepo) Update(o *Order) error {
	if _, ok := r.byID[o.ID]; !ok {
		return ErrNotFound
	}
	c := *o
	r.byID[o.ID] = &c
	return nil
}

func (r *memoryRepo) ListMine(userID string) ([]Order, error) {
	out := []Order{}
	for _, o := range r.byID {
		if o.ClientID == userID || o.ProviderID == userID {
			out = append(out, *o)
		}
	}
	return out, nil
}
