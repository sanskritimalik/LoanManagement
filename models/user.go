package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
User - User consists of all the other essential details associated with a user.
It maintains consistency with Account document on ID.
*/
type User struct {
	ID      primitive.ObjectID `bson:"_id"`
	Email   string             `bson:"email" json:"email"`
	Name    string             `bson:"name" json:"name"`
	Age     int                `bson:"age" json:"age"`
	IsAdmin bool               `bson:"is_admin" json:"is_admin"`
}
