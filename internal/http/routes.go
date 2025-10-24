package http

import (
	"net/http"

	userhttp "github.com/Gab-Mello/service-finder/internal/http/user"
	"github.com/Gab-Mello/service-finder/internal/user"
)

func RegisterAll(mux *http.ServeMux, userSvc *user.Service) {
	// health
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// user
	uh := userhttp.NewHandler(userSvc)
	userhttp.Register(mux, uh)
}
