package controllers

import (
	"bytes"
	"encoding/json"

	"go-gorilla-mongo/cmd/api/configs"
	"go-gorilla-mongo/cmd/api/models"
	"go-gorilla-mongo/cmd/api/schema"
	"go-gorilla-mongo/cmd/api/utils"

	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

}

type RegisterRequest struct {
	User     json.RawMessage `json:"user"`
	UserType string          `json:"user_type"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	userCollection := configs.GetCollection(configs.DB, "users")
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	var user schema.User
	err = json.NewDecoder(bytes.NewReader(request.User)).Decode(&user)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	user.Password = hashedPassword
	result, err := models.Create(userCollection, user)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, result, "data")
}
