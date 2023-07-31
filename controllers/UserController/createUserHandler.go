package controllers

import (
	"aspire/loanmanagement/configs/database"
	"aspire/loanmanagement/configs/log"
	"aspire/loanmanagement/models"
	"aspire/loanmanagement/pkg"
	"aspire/loanmanagement/responses"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser - Takes all attributes associated with a user and adds an
// entry to User document
func CreateUser(c echo.Context) error {
	user := &models.User{}
	if err := pkg.ReadRequestBody(c, user); err != nil {
		msg := "failed to read request body"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusBadRequest, msg)
	}

	user.ID = c.Get("AccountId").(primitive.ObjectID)
	user.IsAdmin = false
	_, err := database.Collections.Users.InsertOne(context.Background(), user)
	if err != nil {
		msg := "failed to post a user"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	}
	return c.JSON(http.StatusCreated, user)
}
