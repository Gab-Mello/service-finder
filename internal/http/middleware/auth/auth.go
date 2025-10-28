package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Gab-Mello/service-finder/internal/auth"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func UserIDFromContext(r *http.Request) (string, bool) {
	uid, ok := r.Context().Value(userIDKey).(string)
	return uid, ok
}

func WithAuth(sessions *auth.SessionManager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sid")
		if err != nil {
			writeJSONErr(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		uid, ok := sessions.Get(c.Value)
		if !ok {
			writeJSONErr(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next(w, r.WithContext(ctx))
	}
}

func writeJSONErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
