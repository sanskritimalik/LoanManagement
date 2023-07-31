package controllers

import (
	"aspire/loanmanagement/configs/log"
	"aspire/loanmanagement/pkg"
	"aspire/loanmanagement/responses"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApproveLoan - Approve loan associated with the given loanId.
// Only Admin users have access to this endpoint
func ApproveLoan(c echo.Context) error {
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
		err = updateLoanByLoanId(c)
		if err != nil {
			msg := "error in updating loan for given loan Id"
			log.Logger.Errorf("%s: %w", msg, err)
			return responses.Message(c, http.StatusInternalServerError, msg)
		}
	}
	return c.JSON(http.StatusOK, "loan status updated successfully")
}

func updateLoanByLoanId(c echo.Context) error {
	loanID, err := primitive.ObjectIDFromHex(c.Param("loanId"))
	if err != nil {
		return fmt.Errorf("Error converting string to ObjectID: %w", err)
	}

	filter := bson.M{"_id": bson.M{"$eq": loanID}}
	update := bson.M{"$set": bson.M{"status": "APPROVED"}}
	return updateLoan(context.Background(), filter, update)
}
