package usermodel

import (
	"context"
	"fmt"
	"go-gorilla-mongo/cmd/api/configs"
	models "go-gorilla-mongo/cmd/api/models/schema"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func Create(_doc interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := userCollection.InsertOne(ctx, _doc)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	return res, nil
}

func FindOne(filter interface{}, options *options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := userCollection.FindOne(ctx, filter, options)
	return res
}

func Find(filter interface{}, options *options.FindOptions) []models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, curError := userCollection.Find(ctx, filter, options)
	if curError != nil {
		panic(curError)
	}
	defer cursor.Close(ctx)
	var users []models.User
	err := cursor.All(ctx, &users)
	if err != nil {
		panic(err)
	}
	return users
}
