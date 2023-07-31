package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Repayment - a repayment is an installment associated with a loan.
*/
type Repayment struct {
	LoanID      primitive.ObjectID `bson:"_loan_id" json:"loan_id"`
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Amount      int                `bson:"amount" json:"amount"`
	ScheduledAt time.Time          `bson:"schedule" json:"schedule"`
	Status      string             `bson:"status" json:"status"`
}
