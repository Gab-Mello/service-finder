package review

import (
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("review not found")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidFields  = errors.New("invalid fields")
	ErrAlreadyExists  = errors.New("review already exists for this order")
	ErrEditWindowOver = errors.New("edit window exceeded")
	ErrOrderNotDone   = errors.New("order not completed")
)

type Review struct {
	OrderID    string    `json:"orderId"`
	ClientID   string    `json:"clientId"`
	ProviderID string    `json:"providerId"`
	Stars      int       `json:"stars"` // 1..5
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
