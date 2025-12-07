package posting

import "time"

type Posting struct {
	ID           string    `json:"ID"`
	ProviderID   string    `json:"ProviderID"`
	ProviderName string    `json:"ProviderName"`
	Title        string    `json:"Title"`
	Description  string    `json:"Description"`
	Price        int64     `json:"Price"`
	Category     string    `json:"Category"`
	City         string    `json:"City"`
	District     string    `json:"District"`
	Archived     bool      `json:"Archived"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
	ProviderAvg  float64   `json:"ProviderAvg,omitempty"`
}

var (
	ErrNotFound      = errStr("posting not found")
	ErrForbidden     = errStr("forbidden")
	ErrInvalidFields = errStr("missing required fields")
)

type errStr string

func (e errStr) Error() string { return string(e) }
