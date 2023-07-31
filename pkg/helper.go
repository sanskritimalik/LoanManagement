package pkg

import (
	"aspire/loanmanagement/configs/database"
	"aspire/loanmanagement/models"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func ReadRequestBody(c echo.Context, object interface{}) error {
	body, err := io.ReadAll(c.Request().Body)
	switch {
	case err != nil:
		return fmt.Errorf("failed to read request body: %w", err)
	case len(body) == 0:
		return fmt.Errorf("request body is empty: %w", echo.ErrBadRequest)
	default:
		if err = json.Unmarshal(body, object); err != nil {
			return fmt.Errorf("failed to unmarshal object: %w", err)
		}
		return nil
	}
}

func IsAdmin(c echo.Context) (bool, error) {
	AccountId := c.Get("AccountId")
	account := &models.Account{}
	if err := database.Collections.Accounts.FindOne(context.Background(), bson.M{"_id": AccountId}).Decode(&account); err != nil {
		panic(err)
	}
	return account.IsAdmin, nil
}
