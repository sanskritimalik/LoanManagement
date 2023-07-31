package environment

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func Init() error {
	log.Print("Initializing environment variables")

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}

	return nil
}
