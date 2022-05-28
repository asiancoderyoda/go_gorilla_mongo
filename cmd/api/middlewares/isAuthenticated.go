package middlewares

import (
	"fmt"
	"go-gorilla-mongo/cmd/api/configs"
	"go-gorilla-mongo/cmd/api/utils"
	"log"
	"net/http"
	"strings"
	"time"

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
		isValidAccessToken, newAccessToken, newRefreshToken, err := VerifyTokens(accessToken, refreshToken)
		if err != nil {
			log.Println(err)
			utils.WriteError(w, err)
		}
		if !isValidAccessToken && newAccessToken == "" && newRefreshToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}
		if !isValidAccessToken && newAccessToken != "" && newRefreshToken != "" {
			w.Header().Set("Authorization", "Bearer "+newAccessToken+":"+newRefreshToken)
		}
		next.ServeHTTP(w, r)
	})
}

func VerifyTokens(accessTokenString string, refreshTokenString string) (bool, string, string, error) {
	accessToken, err := jwt.Parse(accessTokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", parsedToken.Header["alg"])
		}
		return []byte(configs.GetEnvFromKey("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return false, "", "", err
	}
	refreshToken, err := jwt.Parse(refreshTokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", parsedToken.Header["alg"])
		}
		return []byte(configs.GetEnvFromKey("REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return false, "", "", err
	}

	if !accessToken.Valid || !refreshToken.Valid {
		return false, "", "", nil
	}

	accessExpiryTime := accessToken.Claims.(jwt.MapClaims)["exp"].(float64)
	refreshExpiryTime := refreshToken.Claims.(jwt.MapClaims)["exp"].(float64)
	if accessExpiryTime < float64(time.Now().Unix()) && refreshExpiryTime < float64(time.Now().Unix()) {
		return false, "", "", nil
	}
	if accessExpiryTime < float64(time.Now().Unix()) && refreshExpiryTime > float64(time.Now().Unix()) {
		userId := refreshToken.Claims.(jwt.MapClaims)["user"].(string)
		newAccessToken, newRefreshToken, err := utils.GenerateAuthToken(userId)
		if err != nil {
			return false, "", "", err
		}
		return true, newAccessToken, newRefreshToken, nil
	}

	return true, "", "", nil
}
