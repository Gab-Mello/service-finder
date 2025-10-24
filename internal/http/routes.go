package http

import (
	"net/http"

	postinghttp "github.com/Gab-Mello/service-finder/internal/http/posting"
	userhttp "github.com/Gab-Mello/service-finder/internal/http/user"

	"github.com/Gab-Mello/service-finder/internal/posting"
	"github.com/Gab-Mello/service-finder/internal/user"
)

func RegisterAll(mux *http.ServeMux, userSvc *user.Service, postingSvc *posting.Service) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	userHandler := userhttp.NewHandler(userSvc)
	userhttp.Register(mux, userHandler)

	postingHandler := postinghttp.NewHandler(postingSvc)
	postinghttp.Register(mux, postingHandler)
}
