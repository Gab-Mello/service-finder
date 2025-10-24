package user

import "time"

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}

type Service struct {
	repo  Repository
	pw    PasswordHasher
	now   func() time.Time
	idgen func() string
}

func NewService(repo Repository, hasher PasswordHasher, now func() time.Time, idgen func() string) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	if idgen == nil {
		idgen = func() string { return "" }
	}
	return &Service{
		repo:  repo,
		pw:    hasher,
		now:   now,
		idgen: idgen,
	}
}
