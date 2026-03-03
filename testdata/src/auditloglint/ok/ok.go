// Package ok contains valid log calls that should not trigger any diagnostics.
package ok

import (
	"log/slog"

	"go.uber.org/zap"
)

func slogValid() {
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Debug("request received")
	slog.Warn("something went wrong")
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("token validated")
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("retrying connection")
	slog.Info("processing item", "id", 42)
	slog.Info("done")
}

func zapValid() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck
	logger.Info("starting server")
	logger.Error("failed to connect")
	logger.Debug("debug info")
	logger.Warn("resource limit approaching")
}

func zapSugaredValid() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	sugar.Infow("starting server", "port", 8080)
	sugar.Errorw("failed to connect", "err", "timeout")
	sugar.Infof("listening on port %d", 8080)
}
