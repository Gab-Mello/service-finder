package user

import "time"

type Role string

const (
	RoleProvider Role = "provider"
	RoleCustomer Role = "customer"
)

type User struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Email        string           `json:"email"`
	PasswordHash string           `json:"-"`
	Role         Role             `json:"role"`
	Provider     *ProviderProfile `json:"provider,omitempty"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

type ProviderProfile struct {
	Bio       string `json:"bio,omitempty"`
	Phone     string `json:"phone"`
	Expertise string `json:"expertise,omitempty"`
	City      string `json:"city"`
	District  string `json:"district"`
}

var (
	ErrEmailTaken   = errStr("email already in use")
	ErrNotFound     = errStr("user not found")
	ErrUnauthorized = errStr("unauthorized")
)

type errStr string

func (e errStr) Error() string { return string(e) }
