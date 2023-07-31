package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/* Account - stores the login related credentials of the user in DB.
An entry is added to the Accounts collection when the user tries to
access any endpoint after sign up */
type Account struct {
	ID primitive.ObjectID `bson:"_id"`

	Email string `bson:"email"`

	IsAdmin bool `bson:"is_admin"`
}
