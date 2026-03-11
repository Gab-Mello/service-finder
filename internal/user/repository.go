package user

import "sync"

type Repository interface {
	Create(u *User) error
	ByEmail(email string) (*User, error)
	ByID(id string) (*User, error)
	Update(u *User) error
}

type memoryRepo struct {
	mu      sync.RWMutex
	byID    map[string]User
	byEmail map[string]string // email -> id index for O(1) lookup
}

func NewRepository() Repository {
	return &memoryRepo{
		byID:    make(map[string]User),
		byEmail: make(map[string]string),
	}
}

func (r *memoryRepo) Create(u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[u.Email]; exists {
		return ErrEmailTaken
	}

	r.byID[u.ID] = *u
	r.byEmail[u.Email] = u.ID
	return nil
}

func (r *memoryRepo) ByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[email]
	if !ok {
		return nil, ErrNotFound
	}
	it := r.byID[id]
	return &it, nil
}

func (r *memoryRepo) ByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	it, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &it, nil
}

func (r *memoryRepo) Update(u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	old, ok := r.byID[u.ID]
	if !ok {
		return ErrNotFound
	}

	if old.Email != u.Email {
		if _, exists := r.byEmail[u.Email]; exists {
			return ErrEmailTaken
		}
		delete(r.byEmail, old.Email)
		r.byEmail[u.Email] = u.ID
	}

	r.byID[u.ID] = *u
	return nil
}
