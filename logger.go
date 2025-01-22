package panicrecovery

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		// Fall back to nop logger if production logger fials
		Logger = zap.NewNop()
	}
}
