package controllers

import (
	"context"
	"fmt"
	"time"

	"aspire/loanmanagement/configs/database"
	"aspire/loanmanagement/configs/log"
	"aspire/loanmanagement/models"
	"aspire/loanmanagement/pkg"
	"aspire/loanmanagement/responses"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type loanRequest struct {
	Amount int `bson:"amount" json:"amount"`
	Terms  int `bson:"Terms" json:"Terms"`
}

// CreateLoan - Creates a loan corresponding to the logged in user
func CreateLoan(c echo.Context) error {
	AccountId := c.Get("AccountId")
	req := &loanRequest{}
	if err := pkg.ReadRequestBody(c, req); err != nil {
		msg := "failed to read request body"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusBadRequest, msg)
	}

	loan := &models.Loan{}
	loan.SetLoan(primitive.NewObjectID(), req.Amount, req.Terms, "PENDING", AccountId.(primitive.ObjectID))
	repayments, err := generateScheduledRepayments(loan)
	if err != nil {
		msg := "failed to generate repayments"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	}
	loan.AddRepaymentsForLoan(repayments)

	err = createLoanInTransactionalBlock(context.Background(), loan, repayments)
	if err != nil {
		msg := "failed to create a loan and repayments"
		log.Logger.Errorf("%v: %v", msg, err)
		return responses.Message(c, http.StatusInternalServerError, msg)
	}
	return c.JSON(http.StatusCreated, loan)
}

func createLoanInTransactionalBlock(ctx context.Context, loan *models.Loan, repayments []*models.Repayment) error {
	// Start the transaction.
	var repaymentList []interface{}
	for _, repayment := range repayments {
		repaymentList = append(repaymentList, repayment)
	}
	session, err := database.InitializeWithDBSession()
	if err != nil {
		fmt.Println("Error starting session:", err)
		return err
	}
	defer session.EndSession(ctx)

	childCtx := mongo.NewSessionContext(ctx, session)
	_, err = session.WithTransaction(childCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Perform your multiple database updates within this block.

		_, err := database.Collections.Loans.InsertOne(sessCtx, loan)
		if err != nil {
			msg := "failed to post a loan"
			log.Logger.Errorf("%v: %v", msg, err)
			return nil, err
		}

		_, err = database.Collections.Repayments.InsertMany(sessCtx, repaymentList)
		if err != nil {
			msg := "failed to post repayments"
			log.Logger.Errorf("%v: %v", msg, err)
			return nil, err
		}

		// If all updates succeed, the transaction will be committed automatically.

		return nil, nil
	}, nil)

	if err != nil {
		fmt.Println("Error executing transaction:", err)
		return err
	}

	fmt.Println("Transaction executed successfully.")
	return nil
}

func generateScheduledRepayments(loan *models.Loan) ([]*models.Repayment, error) {
	repayments := []*models.Repayment{}
	dueDate := time.Now().AddDate(0, 0, 7)
	for i := 1; i <= loan.Terms; i++ {
		var repayment = &models.Repayment{}
		repayment.ID = primitive.NewObjectID()
		repayment.LoanID = loan.ID
		repayment.Amount = int(loan.Amount) / int(loan.Terms)
		repayment.ScheduledAt = dueDate
		repayment.Status = "PENDING"
		dueDate = dueDate.AddDate(0, 0, 7)
		repayments = append(repayments, repayment)
	}
	return repayments, nil
}
