package middleware

import (
	"context"
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
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		uid, ok := sessions.Get(c.Value)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next(w, r.WithContext(ctx))
	}
}
