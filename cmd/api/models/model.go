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

func Create(ctx context.Context, collection *mongo.Collection, _doc interface{}) (*mongo.InsertOneResult, error) {
	res, err := collection.InsertOne(ctx, _doc)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	return res, nil
}

func FindOne(ctx context.Context, collection *mongo.Collection, filter interface{}, options *options.FindOneOptions) *mongo.SingleResult {
	res := collection.FindOne(ctx, filter, options)
	return res
}

func Find(ctx context.Context, collection *mongo.Collection, filter interface{}, options *options.FindOptions) *mongo.Cursor {
	cursor, curError := collection.Find(ctx, filter, options)
	if curError != nil {
		panic(curError)
	}
	return cursor
	// var users []models.User
	// err := cursor.All(ctx, &users)
	// if err != nil {
	// 	panic(err)
	// }
	// return users
}
