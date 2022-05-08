package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

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
	Token string      `json:"token"`
	User  schema.User `json:"user"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	// userCollection := configs.GetCollection(configs.DB, "users")
	var params LoginRequest
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
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
			"_id": 0,
		},
	}
	// filter := map[string]interface{}{
	// 	"email": user.Email,
	// }
	filter := bson.M{"email": user.Email}
	userCursor := models.Find(userCollection, filter, findOptions)
	err = userCursor.All(context.TODO(), &existingUsers)
	if err != nil {
		utils.WriteError(w, err)
	}
	fmt.Println(existingUsers)
	if len(existingUsers) > 0 {
		fmt.Println(existingUsers)
		utils.WriteError(w, fmt.Errorf("user already exists"))
		return
	}

	user.Password = hashedPassword
	user.ID = models.GenerateID()
	user.CreatedAt = models.GenerateTimestamp()
	user.UpdatedAt = models.GenerateTimestamp()
	_, err = models.Create(userCollection, user)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, user.ID.Hex(), "user")
}
