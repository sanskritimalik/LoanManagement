package controllers

import (
	"aspire/loanmanagement/configs/database"
	"aspire/loanmanagement/configs/log"
	"aspire/loanmanagement/models"
	"aspire/loanmanagement/responses"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// GetCustomerLoans - Fetch all loans associated with the logged in user
func GetCustomerLoans(c echo.Context) error {
	AccountID := c.Get("AccountId")
	loans := []models.Loan{}
	LoansFromDb, err := database.Collections.Loans.Find(context.Background(), bson.M{"user_id": AccountID})
	if err != nil {
		msg := "failed to list users"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	}

	if err := LoansFromDb.All(context.Background(), &loans); err != nil {
		msg := "failed to list users"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	}

	return c.JSON(http.StatusOK, loans)
}
