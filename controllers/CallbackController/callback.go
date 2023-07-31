package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	auth0Config "aspire/loanmanagement/configs/auth0"
	"aspire/loanmanagement/responses"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error", err)
	}
	log.Println("loaded")
}

func FetchJWT(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return responses.Message(c, http.StatusBadRequest, "code parameter was not provided")
	}

	url := auth0Config.TokenFetchURL
	data := auth0Config.GetDataForTokenFetchWithCode(code)

	response, err := http.PostForm(url, data)
	if err != nil {
		return fmt.Errorf("failed to retrieve JWT token from Auth0 server: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response form Auth0 server: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusForbidden {
			return responses.Message(c, http.StatusForbidden, "got response from auth0: unauthorized")
		} else {
			return fmt.Errorf("got bad response from auth0: %v", err)
		}
	}

	fieldsToCheck := struct {
		Scope string `json:"scope"`
	}{}

	if err := json.Unmarshal(body, &fieldsToCheck); err != nil {
		return fmt.Errorf("failed to unmarshal body for field check: %v", err)
	}

	if !strings.Contains(fieldsToCheck.Scope, "email") {
		return responses.Message(c, http.StatusBadRequest, `"email" scope is required`)
	}

	return c.String(http.StatusOK, string(body))
}
