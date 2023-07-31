package responses

import "github.com/labstack/echo/v4"

// returns JSON with message field
func Message(c echo.Context, code int, msg string) error {
	message := struct {
		Message string `json:"message"`
	}{msg}
	return c.JSON(code, message)
}
