package user

type Repository interface {
	Create(u *User) error
	ByEmail(email string) (*User, error)
	ByID(id string) (*User, error)
	Update(u *User) error
}

type memoryRepo struct {
	byID map[string]User
}

func NewRepository() Repository {
	return &memoryRepo{byID: make(map[string]User)}
}

func (r *memoryRepo) Create(u *User) error {
	for _, it := range r.byID {
		if it.Email == u.Email {
			return ErrEmailTaken
		}
	}
	r.byID[u.ID] = *u
	return nil
}

func (r *memoryRepo) ByEmail(email string) (*User, error) {
	for _, it := range r.byID {
		if it.Email == email {
			uu := it
			return &uu, nil
		}
	}
	return nil, ErrNotFound
}

func (r *memoryRepo) ByID(id string) (*User, error) {
	it, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	uu := it
	return &uu, nil
}

func (r *memoryRepo) Update(u *User) error {
	if _, ok := r.byID[u.ID]; !ok {
		return ErrNotFound
	}
	r.byID[u.ID] = *u
	return nil
}
