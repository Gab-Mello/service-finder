package posting

import "time"

type Posting struct {
	ID           string
	ProviderID   string
	ProviderName string
	Title        string
	Description  string
	Price        int64
	Category     string
	City         string
	District     string
	Archived     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var (
	ErrNotFound      = errStr("posting not found")
	ErrForbidden     = errStr("forbidden")
	ErrInvalidFields = errStr("missing required fields")
)

type errStr string

func (e errStr) Error() string { return string(e) }
