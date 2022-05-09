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
		tokens := strings.Split(splittedHeader[1], ":")
		accessToken := tokens[0]
		refreshToken := tokens[1]
		if accessToken == "" || refreshToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
