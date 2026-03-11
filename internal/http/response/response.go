package response

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error encoding JSON response: %v", err)
	}
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

func InternalError(w http.ResponseWriter, err error) {
	log.Printf("internal error: %v", err)
	Error(w, http.StatusInternalServerError, "internal server error")
}

func PathParam(path, prefix, suffix string) string {
	path = strings.TrimPrefix(path, prefix)
	if suffix != "" {
		path = strings.TrimSuffix(path, suffix)
	}
	return path
}

type ErrorMapper func(error) int

func MapError(mappings map[error]int) ErrorMapper {
	return func(err error) int {
		for e, status := range mappings {
			if err == e {
				return status
			}
		}
		return http.StatusInternalServerError
	}
}
