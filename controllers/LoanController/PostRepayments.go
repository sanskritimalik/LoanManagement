package controllers

import (
	"aspire/loanmanagement/configs/database"
	"aspire/loanmanagement/configs/log"
	"aspire/loanmanagement/models"
	"aspire/loanmanagement/pkg"
	"aspire/loanmanagement/responses"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type repaymentRequest struct {
	Amount int `bson:"amount" json:"amount"`
}

// PostRepayment - Pay an installment associated with a given loanId
func PostRepayment(c echo.Context) error {
	loan := &models.Loan{}

	//find the loan corresponding to a the loan id for repayment
	loanID, err := primitive.ObjectIDFromHex(c.Param("loanId"))
	if err != nil {
		msg := "failed to convert loanId into object"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	}

	err = database.Collections.Loans.FindOne(c.Request().Context(), bson.M{"_id": loanID}).Decode(loan)
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		msg := "no loan corresponding to the given id"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusNotFound, msg)
	case err != nil:
		msg := "error fetching loan details"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	default:
		log.Logger.Info("loan details fetched successfully %s", loan.ID)
	}

	AccountId := c.Get("AccountId")
	if loan.UserID != AccountId {
		msg := "user cannot pay loan for given loan id"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusForbidden, msg)
	}

	repaymentReq := &repaymentRequest{}
	if err := pkg.ReadRequestBody(c, repaymentReq); err != nil {
		msg := "failed to read request body"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusBadRequest, msg)
	}

	var pendingRepayment *models.Repayment
	for _, repayment := range loan.Repayments {
		if repayment.Status == "PENDING" {
			pendingRepayment = repayment
			break
		}
	}
	if pendingRepayment == nil {
		msg := "no pending repayments found for this loan"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusMethodNotAllowed, msg)
	}

	// if repaymentReq.Amount.Cmp(pendingRepayment.Amount) >= 0 {
	if repaymentReq.Amount >= pendingRepayment.Amount {
		pendingRepayment.Status = "PAID"
		allPaid := true
		for _, repayment := range loan.Repayments {
			if repayment.Status != "PAID" {
				if pendingRepayment.ID == repayment.ID {
					repayment.Status = "PAID"
					continue
				}
				allPaid = false
				break
			}
		}

		switch {
		case !allPaid:
			repaymentFilter := bson.M{"_id": bson.M{"$eq": pendingRepayment.ID}}
			updatedRepayment := bson.M{"$set": bson.M{"status": "PAID"}}

			loanFilter := bson.M{"_id": bson.M{"$eq": loan.ID}}
			updatedLoan := bson.M{"$set": bson.M{"repayments": loan.Repayments}}
			err = updateRepaymentAndLoanInTransaction(context.Background(), loanFilter, updatedLoan, repaymentFilter, updatedRepayment)
			if err != nil {
				msg := "failed to update repayment and loan"
				log.Logger.Errorf("%v: %v", msg, err)
				return responses.Message(c, http.StatusInternalServerError, msg)
			}
		default:
			repaymentFilter := bson.M{"_id": bson.M{"$eq": pendingRepayment.ID}}
			updatedRepayment := bson.M{"$set": bson.M{"status": "PAID"}}

			loanFilter := bson.M{"_id": bson.M{"$eq": loan.ID}}
			updatedLoan := bson.M{"$set": bson.M{"repayments": loan.Repayments,
				"status": "PAID"}}
			err = updateRepaymentAndLoanInTransaction(context.Background(), loanFilter, updatedLoan, repaymentFilter, updatedRepayment)
			if err != nil {
				msg := "failed to update repayment and loan"
				log.Logger.Errorf("%v: %v", msg, err)
				return responses.Message(c, http.StatusInternalServerError, msg)
			}
		}
		return c.JSON(http.StatusOK, "repayment paid successfully")
	}
	return c.JSON(http.StatusPreconditionFailed, "input amount is less than repayment amount")
}

func updateRepayment(ctx context.Context, repaymentFilter interface{}, repaymentUpdate interface{}) error {
	_, err := database.Collections.Repayments.UpdateOne(ctx, repaymentFilter, repaymentUpdate)
	if err != nil {
		msg := "failed to update repayment"
		log.Logger.Errorf("%v: %v", msg, err)
		return err
	}
	return nil
}

func updateLoan(ctx context.Context, loanFilter interface{}, loanUpdate interface{}) error {
	_, err := database.Collections.Loans.UpdateOne(ctx, loanFilter, loanUpdate)
	if err != nil {
		msg := "failed to update loan"
		log.Logger.Errorf("%v: %v", msg, err)
		return err
	}
	return nil
}

func updateRepaymentAndLoanInTransaction(ctx context.Context, loanFilter interface{}, updatedLoan interface{}, repaymentFilter interface{}, updatedRepayment interface{}) error {
	// Define the options for the transaction.
	opts := options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	session, err := database.InitializeWithDBSession()
	if err != nil {
		return fmt.Errorf("Error starting session: %w", err)
	}
	defer session.EndSession(ctx)

	childCtx := mongo.NewSessionContext(ctx, session)
	_, err = session.WithTransaction(childCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Perform multiple database updates within this block.
		if err := updateRepayment(sessCtx, repaymentFilter, updatedRepayment); err != nil {
			return nil, err
		}

		if err := updateLoan(sessCtx, loanFilter, updatedLoan); err != nil {
			return nil, err
		}

		// If all updates succeed, the transaction will be committed automatically.
		return nil, nil
	}, opts)

	if err != nil {
		return fmt.Errorf("Error executing transaction: %w", err)
	}
	return nil
}
