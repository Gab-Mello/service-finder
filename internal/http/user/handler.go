package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Gab-Mello/service-finder/internal/auth"
	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
	"github.com/Gab-Mello/service-finder/internal/http/response"
	domain "github.com/Gab-Mello/service-finder/internal/user"
)

type Handler struct {
	svc      *domain.Service
	sessions *auth.SessionManager
}

func NewHandler(svc *domain.Service, sessions *auth.SessionManager) *Handler {
	return &Handler{svc: svc, sessions: sessions}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.svc.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		mapErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{
		"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.svc.Authenticate(req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	sid, err := h.sessions.New(u.ID)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	h.sessions.SetCookie(w, sid)

	response.JSON(w, http.StatusOK, map[string]any{
		"userId": u.ID, "name": u.Name, "role": u.Role,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("sid"); err == nil {
		h.sessions.Delete(c.Value)
	}
	h.sessions.ClearCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
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
	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateProviderProfile(w http.ResponseWriter, r *http.Request) {
	uid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ProviderProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.svc.UpdateProviderProfile(uid, domain.ProviderProfile{
		Bio: req.Bio, Phone: req.Phone, Expertise: req.Expertise, City: req.City, District: req.District,
	})
	if err != nil {
		mapErr(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"status": "ok", "provider": u.Provider})
}

func mapErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		response.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrValidation):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.InternalError(w, err)
	}
}
