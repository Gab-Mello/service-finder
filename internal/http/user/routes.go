package user

import (
	"net/http"

	"github.com/Gab-Mello/service-finder/internal/auth"
	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
)

func Register(mux *http.ServeMux, h *Handler, sessions *auth.SessionManager) {
	const api = "/api/v1"
	mux.HandleFunc("POST "+api+"/users", h.Register)
	mux.HandleFunc("POST "+api+"/login", h.Login)
	mux.HandleFunc("POST "+api+"/logout", h.Logout)
	mux.HandleFunc("GET "+api+"/me", authmw.WithAuth(sessions, h.Me))
	mux.HandleFunc("PATCH "+api+"/providers/profile", authmw.WithAuth(sessions, h.UpdateProviderProfile))
}
