package user

import "net/http"

func Register(mux *http.ServeMux, h *Handler) {
	const api = "/api/v1"
	mux.HandleFunc("POST "+api+"/users", h.Register)
	mux.HandleFunc("POST "+api+"/login", h.Login)
	mux.HandleFunc("POST "+api+"/logout", h.Logout)
	mux.HandleFunc("GET "+api+"/me", h.Me)
	mux.HandleFunc("PATCH "+api+"/providers/profile", h.UpdateProviderProfile)
}
