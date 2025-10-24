package user

import (
	"encoding/json"
	"net/http"
	"strings"

	domain "github.com/Gab-Mello/service-finder/internal/user"
)

type Handler struct{ svc *domain.Service }

func NewHandler(svc *domain.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	u, err := h.svc.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		mapErr(w, err)
		return
	}
	writeJSON(w, 201, map[string]any{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	u, err := h.svc.Authenticate(req.Email, req.Password)
	if err != nil {
		writeErr(w, 401, "invalid email or password")
		return
	}
	// sem sessão: devolve o userId para ser usado via ?userId=...
	writeJSON(w, 200, map[string]any{"status": "ok", "userId": u.ID})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"}) // noop por enquanto
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("userId"))
	if id == "" {
		writeErr(w, 400, "missing query parameter: userId")
		return
	}
	u, err := h.svc.ByID(id)
	if err != nil {
		mapErr(w, err)
		return
	}
	resp := map[string]any{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role}
	if u.Provider != nil {
		resp["provider"] = u.Provider
	}
	writeJSON(w, 200, resp)
}

func (h *Handler) UpdateProviderProfile(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("userId"))
	if id == "" {
		writeErr(w, 400, "missing query parameter: userId")
		return
	}
	var req ProviderProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	u, err := h.svc.UpdateProviderProfile(id, domain.ProviderProfile{
		Bio: req.Bio, Phone: req.Phone, Expertise: req.Expertise, City: req.City, District: req.District,
	})
	if err != nil {
		mapErr(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok", "provider": u.Provider})
}

// helpers
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
func mapErr(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrEmailTaken:
		writeErr(w, 400, err.Error())
	case domain.ErrUnauthorized:
		writeErr(w, 401, err.Error())
	case domain.ErrNotFound:
		writeErr(w, 404, err.Error())
	default:
		writeErr(w, 500, "internal error")
	}
}
