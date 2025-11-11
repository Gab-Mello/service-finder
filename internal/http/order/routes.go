package order

import (
	"net/http"
	"strings"

	"github.com/Gab-Mello/service-finder/internal/auth"
	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
)

func Register(mux *http.ServeMux, h *Handler, sessions *auth.SessionManager) {
	const api = "/api/v1"

	mux.HandleFunc("POST "+api+"/orders", authmw.WithAuth(sessions, h.Request))
	mux.HandleFunc("GET "+api+"/orders/mine", authmw.WithAuth(sessions, h.ListMine))
	mux.HandleFunc("GET "+api+"/orders/", authmw.WithAuth(sessions, h.Get))

	mux.HandleFunc("POST "+api+"/orders/", authmw.WithAuth(sessions, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/accept"):
			h.Accept(w, r)
			return
		case strings.HasSuffix(r.URL.Path, "/start"):
			h.Start(w, r)
			return
		case strings.HasSuffix(r.URL.Path, "/complete"):
			h.Complete(w, r)
			return
		case strings.HasSuffix(r.URL.Path, "/cancel"):
			h.Cancel(w, r)
			return
		default:
			http.NotFound(w, r)
		}
	}))
}
