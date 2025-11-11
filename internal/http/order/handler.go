package order

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
	domain "github.com/Gab-Mello/service-finder/internal/order"
)

type Handler struct{ svc *domain.Service }

func NewHandler(s *domain.Service) *Handler { return &Handler{svc: s} }

type requestOrder struct {
	PostingID  string `json:"postingId"`
	ProviderID string `json:"providerId"`
}
type acceptReq struct {
	ScheduledAt string `json:"scheduledAt"`
}

func (h *Handler) Request(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req requestOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PostingID == "" || req.ProviderID == "" {
		writeErr(w, 400, "invalid json")
		return
	}
	o, err := h.svc.Request(uid, req.PostingID, req.ProviderID)
	if err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	writeJSON(w, 201, o)
}

func (h *Handler) Accept(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, 401, "unauthorized")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/orders/")
	id = strings.TrimSuffix(id, "/accept")

	var req acceptReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ScheduledAt == "" {
		writeErr(w, 400, "scheduledAt required (RFC3339)")
		return
	}
	when, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		writeErr(w, 400, "invalid scheduledAt")
		return
	}

	o, err := h.svc.Accept(uid, id, when)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, o)
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, 401, "unauthorized")
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/orders/"), "/start")

	o, err := h.svc.Start(uid, id)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, o)
}

func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, 401, "unauthorized")
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/orders/"), "/complete")

	o, err := h.svc.Complete(uid, id)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, o)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, 401, "unauthorized")
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/orders/"), "/cancel")

	o, err := h.svc.Cancel(uid, id)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, o)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/orders/")
	o, err := h.svc.Get(id)
	if err != nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, o)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, 401, "unauthorized")
		return
	}
	list, _ := h.svc.ListMine(uid)
	writeJSON(w, 200, list)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
func statusFor(err error) int {
	switch err {
	case domain.ErrForbidden:
		return 403
	case domain.ErrInvalidState, domain.ErrInvalidFields:
		return 400
	case domain.ErrNotFound:
		return 404
	default:
		return 500
	}
}
