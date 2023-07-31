package configs

import (
	stdlog "log"

	"aspire/loanmanagement/configs/auth0"
	"aspire/loanmanagement/configs/database"
	"aspire/loanmanagement/configs/environment"
	"aspire/loanmanagement/configs/log"
)

func init() {
	stdlog.Print("Initializing configs")

	if err := log.Init(); err != nil {
		stdlog.Fatalf("failed to initialize logger config: %v", err)
	}

	if err := environment.Init(); err != nil {
		stdlog.Fatalf("failed to initialize environment config: %v", err)
	}

	if err := database.Init(); err != nil {
		stdlog.Fatalf("failed to initialize database config: %v", err)
	}

	if err := auth0.Init(); err != nil {
		stdlog.Fatalf("failed to initialize auth0 config: %v", err)
	}
}
