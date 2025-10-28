package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Gab-Mello/service-finder/internal/auth"
	domain "github.com/Gab-Mello/service-finder/internal/user"
)

type Handler struct {
	svc      *domain.Service
	sessions *auth.SessionManager
}

func NewHandler(svc *domain.Service, sessions *auth.SessionManager) *Handler {
	return &Handler{svc: svc, sessions: sessions}
}

// Register
// @Summary  Cadastrar usuário
// @Tags     users
// @Router   /users [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.svc.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role,
	})
}

// Login
// @Summary  Login
// @Tags     users
// @Router   /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.svc.Authenticate(req.Email, req.Password)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	sid := h.sessions.New(u.ID)
	h.sessions.SetCookie(w, sid)

	writeJSON(w, http.StatusOK, map[string]any{
		"userId": u.ID, "name": u.Name, "role": u.Role,
	})
}

// Logout
// @Summary  Logout
// @Tags     users
// @Router   /logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("sid"); err == nil {
		h.sessions.Delete(c.Value)
	}
	h.sessions.ClearCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// Me
// @Summary  Obter dados do usuário logado
// @Tags     users
// @Router   /me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("sid")
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	uid, ok := h.sessions.Get(c.Value)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	u, err := h.svc.ByID(uid)
	if err != nil {
		mapErr(w, err)
		return
	}

	resp := map[string]any{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role}
	if u.Provider != nil {
		resp["provider"] = u.Provider
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdateProviderProfile
// @Summary  Atualizar perfil do prestador
// @Tags     users
// @Router   /providers/profile [patch]
func (h *Handler) UpdateProviderProfile(w http.ResponseWriter, r *http.Request) {
	// exige estar logado
	c, err := r.Cookie("sid")
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	uid, ok := h.sessions.Get(c.Value)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ProviderProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.svc.UpdateProviderProfile(uid, domain.ProviderProfile{
		Bio: req.Bio, Phone: req.Phone, Expertise: req.Expertise, City: req.City, District: req.District,
	})
	if err != nil {
		mapErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "provider": u.Provider})
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
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		writeErr(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		writeErr(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		writeErr(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrValidation):
		writeErr(w, http.StatusBadRequest, err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal error")
	}
}
