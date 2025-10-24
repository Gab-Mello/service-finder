package posting

import (
	"net/http"
	"strings"
)

func Register(mux *http.ServeMux, h *Handler) {
	const api = "/api/v1"

	mux.HandleFunc("POST "+api+"/postings", h.Create)
	mux.HandleFunc("GET "+api+"/postings", h.ListPublic)
	mux.HandleFunc("GET "+api+"/postings/mine", h.ListMine)

	mux.HandleFunc("GET "+api+"/postings/", h.GetPublic)
	mux.HandleFunc("PATCH "+api+"/postings/", h.Update)
	mux.HandleFunc("POST "+api+"/postings/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/archive") {
			h.Archive(w, r)
			return
		}
		http.NotFound(w, r)
	})
}
