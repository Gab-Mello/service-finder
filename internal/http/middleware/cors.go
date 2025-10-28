package middleware

import "net/http"

func PermissiveCORSWithoutCreds(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Qualquer origem
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Todos os métodos comuns
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// Libera todos os headers requisitados no preflight (ou tudo)
		reqHeaders := r.Header.Get("Access-Control-Request-Headers")
		if reqHeaders == "" {
			w.Header().Set("Access-Control-Allow-Headers", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
		}
		// Cache do preflight
		w.Header().Set("Access-Control-Max-Age", "86400") // 24h
		// Para caches/proxies respeitarem variações de preflight
		w.Header().Set("Vary", "Access-Control-Request-Headers")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Origin")

		// Responde preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
