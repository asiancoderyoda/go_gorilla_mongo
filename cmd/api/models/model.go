package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func GetIdFromHex(_id string) primitive.ObjectID {
	isValid := primitive.IsValidObjectID(_id)
	if !isValid {
		panic("Invalid ObjectID")
	}
	id, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		panic(err)
	}
	return id
}

func GenerateTimestamp() primitive.Timestamp {
	return primitive.Timestamp{
		T: uint32(time.Now().Unix()),
		I: 0,
	}
}

func Create(collection *mongo.Collection, _doc interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, _doc)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	return res, nil
}

func FindOne(collection *mongo.Collection, filter interface{}, options *options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, filter, options)
	return res
}

func Find(collection *mongo.Collection, filter interface{}, options *options.FindOptions) *mongo.Cursor {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, curError := collection.Find(ctx, filter, options)
	if curError != nil {
		panic(curError)
	}
	defer cursor.Close(ctx)
	return cursor
	// var users []models.User
	// err := cursor.All(ctx, &users)
	// if err != nil {
	// 	panic(err)
	// }
	// return users
}