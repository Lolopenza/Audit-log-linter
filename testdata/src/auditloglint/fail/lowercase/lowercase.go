// Package fail_lowercase contains log calls that violate the lowercase rule.
package fail_lowercase

import (
	"log/slog"

	"go.uber.org/zap"
)

func badLowercase() {
	slog.Info("Starting server on port 8080")   // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect to database") // want `log message must start with a lowercase letter`
	slog.Warn("Warning something happened")     // want `log message must start with a lowercase letter`
	slog.Debug("Request received from client")  // want `log message must start with a lowercase letter`

	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck
	logger.Info("Server started")            // want `log message must start with a lowercase letter`
	logger.Error("Connection error occurred") // want `log message must start with a lowercase letter`
}
