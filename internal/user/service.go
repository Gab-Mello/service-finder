package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Gab-Mello/service-finder/internal/ports"
	"github.com/google/uuid"
)

var ErrValidation = errors.New("validation error")

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
		idgen = func() string { return uuid.NewString() }
	}
	if hasher == nil {
		hasher = noOpHasher{}
	}
	return &Service{repo: repo, pw: hasher, now: now, idgen: idgen}
}

func (s *Service) Register(name, email, password, role string) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	role = strings.ToLower(strings.TrimSpace(role))

	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("%w: invalid email", ErrValidation)
	}
	if len(password) < 8 {
		return nil, fmt.Errorf("%w: password must be at least 8 characters", ErrValidation)
	}
	if role != string(RoleProvider) && role != string(RoleCustomer) {
		return nil, fmt.Errorf("%w: role must be 'provider' or 'customer", ErrValidation)
	}

	if _, err := s.repo.ByEmail(email); err == nil {
		return nil, ErrEmailTaken
	}

	hash, err := s.pw.Hash(password)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:           s.idgen(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         Role(role),
		CreatedAt:    s.now(),
		UpdatedAt:    s.now(),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	u, err := s.repo.ByEmail(email)
	if err != nil || !s.pw.Compare(u.PasswordHash, password) {
		return nil, ErrUnauthorized
	}
	return u, nil
}

func (s *Service) UpdateProviderProfile(userID string, p ProviderProfile) (*User, error) {
	u, err := s.repo.ByID(userID)
	if err != nil {
		return nil, err
	}
	if u.Role != RoleProvider {
		return nil, ErrUnauthorized
	}

	if strings.TrimSpace(p.Phone) == "" || strings.TrimSpace(p.City) == "" || strings.TrimSpace(p.District) == "" {
		return nil, errors.New("phone, city and district are required")
	}
	u.Provider = &ProviderProfile{
		Bio: strings.TrimSpace(p.Bio), Phone: strings.TrimSpace(p.Phone),
		Expertise: strings.TrimSpace(p.Expertise), City: strings.TrimSpace(p.City), District: strings.TrimSpace(p.District),
	}
	u.UpdatedAt = s.now()
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

var _ ports.ProviderDirectory = (*Service)(nil)

func (s *Service) GetNameByID(id string) (string, error) {
	u, err := s.repo.ByID(id)
	if err != nil {
		return "", err
	}
	return u.Name, nil
}

func (s *Service) ByID(id string) (*User, error) { return s.repo.ByID(id) }

type noOpHasher struct{}

func (noOpHasher) Hash(p string) (string, error) { return p, nil }
func (noOpHasher) Compare(h, p string) bool      { return h == p }
