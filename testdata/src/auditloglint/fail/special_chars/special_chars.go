// Package fail_special contains log calls that violate the special characters rule.
package fail_special

import (
	"log/slog"

	"go.uber.org/zap"
)

func badSpecialChars() {
	slog.Info("connection failed!!!")        // want `log message must not contain special character`
	slog.Warn("warning: something went wrong") // want `log message must not contain special character`
	slog.Error("something went wrong...")    // want `log message must not contain`

	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck
	logger.Info("server started!")           // want `log message must not contain special character`
}

func badEmoji() {
	slog.Info("server started \U0001F680")   // want `log message must not contain emoji`
	slog.Error("connection failed \U0001F4A5") // want `log message must not contain emoji`
}
