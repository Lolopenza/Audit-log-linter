// Package fail_sensitive contains log calls that violate the sensitive data rule.
package fail_sensitive

import (
	"log/slog"

	"go.uber.org/zap"
)

func badSensitiveConcatenation() {
	password := "secret123"
	apiKey := "key-abc"
	token := "tok-xyz"

	slog.Info("user password: " + password) // want `log message may expose sensitive data`
	slog.Debug("api_key=" + apiKey)         // want `log message may expose sensitive data`
	slog.Info("token: " + token)            // want `log message may expose sensitive data`

	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck
	logger.Info("connecting with " + password) // want `log message may expose sensitive data`
}

func badSensitiveIdentName() {
	password := "s3cr3t"
	slog.Info("processing: " + password) // want `log message may expose sensitive data`

	secret := "topsecret"
	slog.Error("failed: " + secret) // want `log message may expose sensitive data`
}
