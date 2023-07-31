package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Loan - model representation of Loan. Stores all the loan
related information for a user corresponding to the loan id
*/
type Loan struct {
	ID         primitive.ObjectID `bson:"_id"`
	Amount     int                `bson:"amount" json:"amount"`
	Terms      int                `bson:"Terms" json:"Terms"`
	Status     string             `bson:"status" json:"status"`
	UserID     primitive.ObjectID `bson:"user_id"`
	CreatedAt  time.Time          `bson:"creation_time" json:"creation_time"`
	Repayments []*Repayment       `bson:"repayments" json:"repayments"`
}

// SetLoan - sets loan parameters of a Loan object
func (l *Loan) SetLoan(ID primitive.ObjectID,
	Amount int,
	Terms int,
	Status string,
	UserID primitive.ObjectID) {
	l.ID = ID
	l.Amount = Amount
	l.Terms = Terms
	l.CreatedAt = time.Now()
	l.Status = Status
	l.UserID = UserID
	l.Repayments = []*Repayment{}
}

// SetRepaymentsForLoan - adds Repayments corresponding to a Loan
func (l *Loan) AddRepaymentsForLoan(Repayments []*Repayment) {
	l.Repayments = append(l.Repayments, Repayments...)
}
