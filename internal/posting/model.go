package posting

import "time"

type Posting struct {
	ID           string    `json:"id"`
	ProviderID   string    `json:"providerId"`
	ProviderName string    `json:"providerName"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Price        int64     `json:"price"`
	Category     string    `json:"category"`
	City         string    `json:"city"`
	District     string    `json:"district"`
	Archived     bool      `json:"archived"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	ProviderAvg  float64   `json:"providerAvg,omitempty"`
}

var (
	ErrNotFound      = errStr("posting not found")
	ErrForbidden     = errStr("forbidden")
	ErrInvalidFields = errStr("missing required fields")
)

type errStr string

func (e errStr) Error() string { return string(e) }
