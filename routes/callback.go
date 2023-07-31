package routes

import (
	controllers "aspire/loanmanagement/controllers/CallbackController"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error", err)
	}
}

func ApplyCallback(e *echo.Echo) {
	fmt.Println("apply callback")
	fmt.Println(os.Getenv("Auth0CallbackEndpoint"))
	e.GET(os.Getenv("Auth0CallbackEndpoint"), controllers.FetchJWT)
}
