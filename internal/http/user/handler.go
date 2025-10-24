package user

import (
	"encoding/json"
	"net/http"

	domain "github.com/Gab-Mello/service-finder/internal/user"
)

type Handler struct{ svc *domain.Service }

func NewHandler(svc *domain.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(w http.ResponseWriter, r *http.Request)              { notImplemented(w) }
func (h *Handler) Login(w http.ResponseWriter, r *http.Request)                 { notImplemented(w) }
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request)                { notImplemented(w) }
func (h *Handler) Me(w http.ResponseWriter, r *http.Request)                    { notImplemented(w) }
func (h *Handler) UpdateProviderProfile(w http.ResponseWriter, r *http.Request) { notImplemented(w) }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func notImplemented(w http.ResponseWriter) {
	writeJSON(w, 501, map[string]string{"error": "not implemented"})
}
