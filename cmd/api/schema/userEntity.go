package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleCompany
	RoleEmployee
)

func (r Role) getRole() string {
	return [...]string{"admin", "company", "employee"}[r-1]
}

type User struct {
	ID             primitive.ObjectID  `json:"_id" bson:"_id"`
	Name           string              `json:"name" bson:"name"`
	Email          string              `json:"email" bson:"email"`
	Password       string              `json:"password" bson:"password"`
	ProfilePicture string              `json:"profile_picture" bson:"profile_picture"`
	Role           string              `json:"role" bson:"role"`
	CreatedAt      primitive.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt      primitive.Timestamp `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID             string              `json:"_id"`
	Name           string              `json:"name"`
	Email          string              `json:"email"`
	ProfilePicture string              `json:"profile_picture"`
	CreatedAt      primitive.Timestamp `json:"created_at"`
	UpdatedAt      primitive.Timestamp `json:"updated_at"`
}
