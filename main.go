package main

import (
	_ "aspire/loanmanagement/configs"
	"aspire/loanmanagement/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error", err)
	}
	log.Println("loaded")
}

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	routes.ApplyCallback(e)
	routes.ApplyLoan(e)
	e.Logger.Fatal(e.Start(":8080"))
}
