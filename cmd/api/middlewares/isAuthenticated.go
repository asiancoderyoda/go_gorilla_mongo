package middlewares

import (
	"net/http"
	"strings"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		splittedHeader := strings.Split(bearerToken, " ")
		if len(splittedHeader) != 2 || splittedHeader[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		// authToken := splittedHeader[1]

		next.ServeHTTP(w, r)
	})
}
