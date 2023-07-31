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
	"go.mongodb.org/mongo-driver/bson"
)

// GetUsers - Fetches all users registered with application.
// Only Admin can access this functionality
func GetUsers(c echo.Context) error {
	isAdminUser, err := pkg.IsAdmin(c)
	switch {
	case err != nil:
		msg := "error in loan approval for the given loanId"
		log.Logger.Errorf("%s: %w", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	case !isAdminUser:
		msg := "User is not allowed access to the given page"
		log.Logger.Errorf("%s: %w", msg, err)
		return responses.Message(c, http.StatusForbidden, msg)
	default:
		users := []*models.User{}
		usersFromDb, err := database.Collections.Users.Find(context.Background(), bson.M{})
		if err != nil {
			msg := "failed to list users"
			log.Logger.Errorf("%v: %v", msg, err)
			return responses.Message(c, http.StatusInternalServerError, msg)
		}

		if err := usersFromDb.All(context.Background(), &users); err != nil {
			msg := "failed to list users"
			log.Logger.Errorf("%v: %v", msg, err)
			return responses.Message(c, http.StatusInternalServerError, msg)
		}
		return c.JSON(http.StatusOK, users)
	}
}
