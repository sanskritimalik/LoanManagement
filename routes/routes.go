package routes

import (
	"aspire/loanmanagement/configs/auth0"
	LoanController "aspire/loanmanagement/controllers/LoanController"
	UserController "aspire/loanmanagement/controllers/UserController"
	d2middlewares "aspire/loanmanagement/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ApplyLoan(e *echo.Echo) {
	// Fetches all users registered with the application. Only admin can access this endpoint
	e.GET("/user", UserController.GetUsers,
		middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte(auth0.Secret)}),
		d2middlewares.AssociateAccountWithRequest,
	)

	// Creates a new user in Database
	e.POST("/user", UserController.CreateUser,
		middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte(auth0.Secret)}),
		d2middlewares.AssociateAccountWithRequest,
	)

	// Creates a loan corresponding to the logged in user
	e.POST("/loans", LoanController.CreateLoan,
		middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte(auth0.Secret)}),
		d2middlewares.AssociateAccountWithRequest,
	)

	// Fetch all loans associated with the logged in user
	e.GET("/loans", LoanController.GetCustomerLoans,
		middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte(auth0.Secret)}),
		d2middlewares.AssociateAccountWithRequest,
	)

	// Pay an installment associated with a given loanId
	e.POST("/repayments/:loanId", LoanController.PostRepayment,
		middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte(auth0.Secret)}),
		d2middlewares.AssociateAccountWithRequest,
	)

	// Approve loan associated with the given loanId.
	// Only Admin users have access to this endpoint
	e.GET("/approve/:loanId", LoanController.ApproveLoan,
		middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte(auth0.Secret)}),
		d2middlewares.AssociateAccountWithRequest,
	)
}
