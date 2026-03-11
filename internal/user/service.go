package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrValidation = errors.New("validation error")
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^[\d\s\-\(\)\+]{8,20}$`)
)

const (
	minPasswordLen = 8
	maxNameLen     = 100
	maxBioLen      = 1000
	maxPhoneLen    = 20
	maxCityLen     = 100
	maxDistrictLen = 100
)

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
	if len(name) > maxNameLen {
		return nil, fmt.Errorf("%w: name is too long", ErrValidation)
	}
	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("%w: invalid email format", ErrValidation)
	}
	if len(password) < minPasswordLen {
		return nil, fmt.Errorf("%w: password must be at least %d characters", ErrValidation, minPasswordLen)
	}
	if role != string(RoleProvider) && role != string(RoleCustomer) {
		return nil, fmt.Errorf("%w: role must be 'provider' or 'customer'", ErrValidation)
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

	phone := strings.TrimSpace(p.Phone)
	city := strings.TrimSpace(p.City)
	district := strings.TrimSpace(p.District)
	bio := strings.TrimSpace(p.Bio)
	expertise := strings.TrimSpace(p.Expertise)

	if phone == "" || city == "" || district == "" {
		return nil, fmt.Errorf("%w: phone, city and district are required", ErrValidation)
	}
	if !phoneRegex.MatchString(phone) {
		return nil, fmt.Errorf("%w: invalid phone format", ErrValidation)
	}
	if len(bio) > maxBioLen {
		return nil, fmt.Errorf("%w: bio is too long", ErrValidation)
	}
	if len(city) > maxCityLen || len(district) > maxDistrictLen {
		return nil, fmt.Errorf("%w: city or district is too long", ErrValidation)
	}

	u.Provider = &ProviderProfile{
		Bio:       bio,
		Phone:     phone,
		Expertise: expertise,
		City:      city,
		District:  district,
	}
	u.UpdatedAt = s.now()
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

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
