package posting

import (
	"net/http"
	"strings"

	"github.com/Gab-Mello/service-finder/internal/auth"
	middleware "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
)

func Register(mux *http.ServeMux, h *Handler, sessions *auth.SessionManager) {
	const api = "/api/v1"

	mux.HandleFunc("GET "+api+"/postings", h.Search)
	mux.HandleFunc("GET "+api+"/postings/", h.GetPublic)

	mux.HandleFunc("POST "+api+"/postings", middleware.WithAuth(sessions, h.Create))
	mux.HandleFunc("GET "+api+"/postings/mine", middleware.WithAuth(sessions, h.ListMine))
	mux.HandleFunc("PATCH "+api+"/postings/", middleware.WithAuth(sessions, h.Update))
	mux.HandleFunc("POST "+api+"/postings/", middleware.WithAuth(sessions, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/archive") {
			h.Archive(w, r)
			return
		}
		http.NotFound(w, r)
	}))
}
