package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gorilla-mongo/cmd/api/configs"
	"go-gorilla-mongo/cmd/api/models"
	"go-gorilla-mongo/cmd/api/schema"
	"go-gorilla-mongo/cmd/api/utils"

	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
* Login request
* @param {LoginRequest} request
* @return {LoginResponse}
 */
type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string              `json:"accessToken"`
	RefreshToken string              `json:"refreshToken"`
	User         schema.UserResponse `json:"user"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := configs.GetCollection(configs.DB, "users")
	var params LoginRequest
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	var user schema.User
	filter := map[string]interface{}{
		"email": params.UserName,
	}
	item := models.FindOne(ctx, userCollection, filter, nil)
	err = item.Decode(&user)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	isValidUser, err := utils.ValidatePassword(params.Password, user.Password)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	if !isValidUser {
		utils.WriteJSON(w, http.StatusUnauthorized, "invalid username or password", "error")
	}
	accessToken, refreshToken, err := utils.GenerateAuthToken(user.ID.Hex())
	if err != nil {
		utils.WriteError(w, err)
	}
	AuthUser := schema.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	loginResponse := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         AuthUser,
	}
	w.Header().Set("Authorization", "Bearer "+accessToken+":"+refreshToken)
	utils.WriteJSON(w, http.StatusOK, loginResponse, "data")
}

/*
* Register a new user
* @param {RegisterRequest} request
* @return {string} userID
 */
type RegisterRequest struct {
	User     json.RawMessage `json:"user"`
	UserType string          `json:"user_type"`
}

type ExistingUsers struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := configs.GetCollection(configs.DB, "users")
	var params RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	var user schema.User
	err = json.NewDecoder(bytes.NewReader(params.User)).Decode(&user)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	var existingUsers []ExistingUsers
	findOptions := &options.FindOptions{
		Projection: map[string]interface{}{
			"_id": 1,
		},
	}
	filter := bson.M{"email": user.Email}
	userCursor := models.Find(ctx, userCollection, filter, findOptions)
	for userCursor.Next(ctx) {
		var existingUser ExistingUsers
		err = userCursor.Decode(&existingUser)
		if err != nil {
			utils.WriteError(w, err)
		}
		existingUsers = append(existingUsers, existingUser)
	}
	if err != nil {
		utils.WriteError(w, err)
	}
	if len(existingUsers) > 0 {
		utils.WriteError(w, fmt.Errorf("user already exists"))
		return
	}

	user.Password = hashedPassword
	user.ID = models.GenerateID()
	user.CreatedAt = models.GenerateTimestamp()
	user.UpdatedAt = models.GenerateTimestamp()
	_, err = models.Create(ctx, userCollection, user)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, user.ID.Hex(), "user")
}
