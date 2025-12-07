package review

import (
	"net/http"

	"github.com/Gab-Mello/service-finder/internal/auth"
	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
)

func Register(mux *http.ServeMux, h *Handler, sessions *auth.SessionManager) {
	const api = "/api/v1"

	mux.HandleFunc("POST "+api+"/reviews", authmw.WithAuth(sessions, h.Create))
	mux.HandleFunc("PATCH "+api+"/reviews/", authmw.WithAuth(sessions, h.Edit))
}
