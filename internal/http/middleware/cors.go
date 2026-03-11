package middleware

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string
	AllowedHeaders []string
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
		},
	}
}

func CORSWithCreds(next http.Handler) http.Handler {
	return CORSWithConfig(DefaultCORSConfig())(next)
}

func CORSWithConfig(cfg CORSConfig) func(http.Handler) http.Handler {
	allowedOrigins := make(map[string]bool)
	for _, o := range cfg.AllowedOrigins {
		allowedOrigins[o] = true
	}

	allowedHeadersStr := strings.Join(cfg.AllowedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if allowedOrigins[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeadersStr)

			w.Header().Set("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
			w.Header().Add("Vary", "Access-Control-Request-Method")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
