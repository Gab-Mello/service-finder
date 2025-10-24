package http

import (
	"log"
	"net/http"
	"time"
)

func NewServer() *http.ServeMux {
	return http.NewServeMux()
}

func Listen(addr string, handler *http.ServeMux) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           withLogging(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
