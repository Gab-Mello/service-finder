package http

import (
	"net/http"

	userhttp "github.com/Gab-Mello/service-finder/internal/http/user"
	"github.com/Gab-Mello/service-finder/internal/posting"

	"github.com/Gab-Mello/service-finder/internal/auth"
	postinghttp "github.com/Gab-Mello/service-finder/internal/http/posting"
	"github.com/Gab-Mello/service-finder/internal/user"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func RegisterAll(mux *http.ServeMux, sessions *auth.SessionManager, userSvc *user.Service, postingSvc *posting.Service) {

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	uh := userhttp.NewHandler(userSvc, sessions)
	userhttp.Register(mux, uh)

	ph := postinghttp.NewHandler(postingSvc)
	postinghttp.Register(mux, ph, sessions)
}
