// Package main Service Finder API.
//
// @title       Service Finder API
// @version     1.0
// @description Plataforma de prestação e busca de serviços.
// @BasePath    /api/v1
package main

import (
	"log"
	"time"

	"github.com/Gab-Mello/service-finder/internal/auth"
	transport "github.com/Gab-Mello/service-finder/internal/http"
	"github.com/Gab-Mello/service-finder/internal/order"
	"github.com/Gab-Mello/service-finder/internal/posting"
	"github.com/Gab-Mello/service-finder/internal/user"

	_ "github.com/Gab-Mello/service-finder/docs"
)

func main() {
	addr := ":8080"

	sessions := auth.NewSessionManager(2 * time.Hour)

	userRepo := user.NewRepository()
	userSvc := user.NewService(userRepo, nil, time.Now, nil)

	postRepo := posting.NewRepository()
	postSvc := posting.NewService(postRepo, userSvc, time.Now, nil)

	orderRepo := order.NewRepository()
	orderSvc := order.NewService(orderRepo, time.Now, nil, nil)

	mux := transport.NewServer()
	transport.RegisterAll(mux, sessions, userSvc, postSvc, orderSvc)

	log.Printf("listening on %s", addr)
	log.Fatal(transport.Listen(addr, mux))
}
