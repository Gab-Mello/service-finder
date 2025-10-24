package user

import "time"

type Role string

const (
	RoleProvider Role = "provider"
	RoleCustomer Role = "customer"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	Provider     *ProviderProfile
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProviderProfile struct {
	Bio       string
	Phone     string
	Expertise string
	City      string
	District  string
}

var (
	ErrEmailTaken   = errStr("email already in use")
	ErrNotFound     = errStr("user not found")
	ErrUnauthorized = errStr("unauthorized")
)

type errStr string

func (e errStr) Error() string { return string(e) }
