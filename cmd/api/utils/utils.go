package utils

import (
	"encoding/json"
	"go-gorilla-mongo/cmd/api/configs"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func WriteJSON(w http.ResponseWriter, statusCode int, payload interface{}, wrap string) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = payload
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)

	return nil
}

func WriteError(w http.ResponseWriter, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	httpError := errorResponse{
		Error: err.Error(),
	}

	writeErr := WriteJSON(w, http.StatusUnprocessableEntity, httpError, "error")

	if writeErr != nil {
		panic(writeErr)
	}
}

func ValidatePassword(password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GenerateAuthToken(userId string) (string, string, error) {
	hmacAccessKeySecret := []byte(configs.GetEnvFromKey("ACCESS_TOKEN_SECRET"))
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userId,
		"iat":  time.Now().Unix(),                         // issued at
		"nbf":  time.Now().Unix(),                         // valid from this time
		"exp":  time.Now().Add(time.Second * 3600).Unix(), // expires in
		"iss":  "go-gorilla-mongo",
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userId,
		"iat":  time.Now().Unix(),                          // issued at
		"nbf":  time.Now().Unix(),                          // valid from this time
		"exp":  time.Now().Add(time.Second * 86400).Unix(), // expires in
		"iss":  "go-gorilla-mongo",
	})
	signedAccessToken, err := accessToken.SignedString(hmacAccessKeySecret)
	if err != nil {
		return "", "", err
	}
	signedRefreshToken, err := refreshToken.SignedString(hmacAccessKeySecret)
	if err != nil {
		return "", "", err
	}
	return signedAccessToken, signedRefreshToken, nil
}
