package posting

type Repository interface {
	Create(*Posting) error
	Update(*Posting) error
	ByID(id string) (*Posting, error)
	ListByProvider(providerID string) ([]Posting, error)
	ListPublic() ([]Posting, error)
}

type memoryRepo struct {
	byID map[string]Posting
}

func NewRepository() Repository {
	return &memoryRepo{byID: make(map[string]Posting)}
}

func (r *memoryRepo) Create(p *Posting) error {
	r.byID[p.ID] = *p
	return nil
}

func (r *memoryRepo) Update(p *Posting) error {
	if _, ok := r.byID[p.ID]; !ok {
		return ErrNotFound
	}
	r.byID[p.ID] = *p
	return nil
}

func (r *memoryRepo) ByID(id string) (*Posting, error) {
	it, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := it
	return &cp, nil
}

func (r *memoryRepo) ListByProvider(pid string) ([]Posting, error) {
	out := make([]Posting, 0)
	for _, it := range r.byID {
		if it.ProviderID == pid {
			out = append(out, it)
		}
	}
	return out, nil
}

func (r *memoryRepo) ListPublic() ([]Posting, error) {
	out := make([]Posting, 0)
	for _, it := range r.byID {
		if !it.Archived {
			out = append(out, it)
		}
	}
	return out, nil
}
