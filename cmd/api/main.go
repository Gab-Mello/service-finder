package main

import (
	"log"

	transport "github.com/Gab-Mello/service-finder/internal/http"
	"github.com/Gab-Mello/service-finder/internal/user"
)

func main() {
	repo := user.NewRepository()
	svc := user.NewService(repo, nil, nil, nil)

	mux := transport.NewServer()
	transport.RegisterAll(mux, svc)

	log.Printf("listening on %s", ":8080")
	if err := transport.Listen(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
