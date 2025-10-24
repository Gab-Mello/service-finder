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
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
