package order

import (
	"errors"
	"time"
)

type Status string

const (
	StatusPending    Status = "PENDENTE"
	StatusAccepted   Status = "ACEITO"
	StatusInProgress Status = "EM_ANDAMENTO"
	StatusCompleted  Status = "CONCLUIDO"
	StatusCanceled   Status = "CANCELADO"
)

var (
	ErrNotFound      = errors.New("order not found")
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidFields = errors.New("invalid fields")
	ErrInvalidState  = errors.New("invalid state transition")
)

type HistoryEntry struct {
	At   time.Time `json:"at"`
	By   string    `json:"by"`
	From Status    `json:"from"`
	To   Status    `json:"to"`
	Note string    `json:"note,omitempty"`
}

type Order struct {
	ID          string         `json:"id"`
	PostingID   string         `json:"postingId"`
	ClientID    string         `json:"clientId"`
	ProviderID  string         `json:"providerId"`
	ScheduledAt *time.Time     `json:"scheduledAt,omitempty"`
	Status      Status         `json:"status"`
	History     []HistoryEntry `json:"history"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
