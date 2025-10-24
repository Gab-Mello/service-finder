package user

import "sync"

type Repository interface {
	Create(u *User) error
	ByEmail(email string) (*User, error)
	ByID(id string) (*User, error)
	Update(u *User) error
}

type memoryRepo struct {
	mu   sync.RWMutex
	byID map[string]*User
	byEM map[string]string
}

func NewRepository() Repository {
	return &memoryRepo{byID: map[string]*User{}, byEM: map[string]string{}}
}

func (r *memoryRepo) Create(u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byEM[u.Email]; ok {
		return ErrEmailTaken
	}
	c := *u
	r.byID[u.ID] = &c
	r.byEM[u.Email] = u.ID
	return nil
}
func (r *memoryRepo) ByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byEM[email]
	if !ok {
		return nil, ErrNotFound
	}
	c := *r.byID[id]
	return &c, nil
}
func (r *memoryRepo) ByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	c := *u
	return &c, nil
}
func (r *memoryRepo) Update(u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[u.ID]; !ok {
		return ErrNotFound
	}
	c := *u
	r.byID[u.ID] = &c
	return nil
}
