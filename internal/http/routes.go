package http

import (
	"net/http"

	postinghttp "github.com/Gab-Mello/service-finder/internal/http/posting"
	userhttp "github.com/Gab-Mello/service-finder/internal/http/user"

	"github.com/Gab-Mello/service-finder/internal/posting"
	"github.com/Gab-Mello/service-finder/internal/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func RegisterAll(mux *http.ServeMux, userSvc *user.Service, postingSvc *posting.Service) {
	// healthcheck
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	// swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// users
	uh := userhttp.NewHandler(userSvc)
	userhttp.Register(mux, uh)

	// postings
	ph := postinghttp.NewHandler(postingSvc)
	postinghttp.Register(mux, ph)
}
