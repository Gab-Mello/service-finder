package main

import (
	"log"
	"time"

	transport "github.com/Gab-Mello/service-finder/internal/http"
	"github.com/Gab-Mello/service-finder/internal/posting"
	"github.com/Gab-Mello/service-finder/internal/user"
)

func main() {
	addr := ":8080"

	userRepo := user.NewRepository()
	userSvc := user.NewService(userRepo, nil, time.Now, nil)

	postRepo := posting.NewRepository()
	postSvc := posting.NewService(postRepo, time.Now, nil)

	mux := transport.NewServer()
	transport.RegisterAll(mux, userSvc, postSvc)

	log.Printf("listening on %s", addr)
	log.Fatal(transport.Listen(addr, mux))
}
