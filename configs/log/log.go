package log

import (
	"fmt"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to build production logger: %v", err)
	}
	Logger = logger.Sugar()
	return nil
}
