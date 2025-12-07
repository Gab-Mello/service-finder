package review

import (
	"encoding/json"
	"net/http"

	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
	domain "github.com/Gab-Mello/service-finder/internal/review"
)

type Handler struct{ svc *domain.Service }

func NewHandler(s *domain.Service) *Handler { return &Handler{svc: s} }

type createReq struct {
	OrderID string `json:"orderId"`
	Stars   int    `json:"stars"`   // 1..5
	Comment string `json:"comment"` // opcional
}
type editReq struct {
	Stars   int    `json:"stars"`
	Comment string `json:"comment"`
}

// POST /reviews
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, 401, "unauthorized")
		return
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderID == "" {
		writeErr(w, 400, "invalid json")
		return
	}
	rv, err := h.svc.Create(uid, req.OrderID, req.Stars, req.Comment)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 201, rv)
}

// PATCH /reviews/{orderId}
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, 401, "unauthorized")
		return
	}

	orderID := r.URL.Path[len("/api/v1/reviews/"):]
	var req editReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	rv, err := h.svc.Edit(uid, orderID, req.Stars, req.Comment)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, rv)
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
	case domain.ErrInvalidFields, domain.ErrAlreadyExists, domain.ErrEditWindowOver, domain.ErrOrderNotDone:
		return 400
	case domain.ErrNotFound:
		return 404
	default:
		return 500
	}
}
