package middlewares

import (
	"fmt"
	"go-gorilla-mongo/cmd/api/configs"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", parsedToken.Header["alg"])
		}
		return []byte(configs.GetEnvFromKey("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, nil
	}

	return true, nil
}
