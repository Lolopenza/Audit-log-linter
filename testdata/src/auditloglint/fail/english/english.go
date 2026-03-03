// Package fail_english contains log calls that violate the English-only rule.
package fail_english

import (
	"log/slog"

	"go.uber.org/zap"
)

func badEnglish() {
	slog.Info("запуск сервера")                         // want `log message must be in English only`
	slog.Error("ошибка подключения к базе данных")      // want `log message must be in English only`
	slog.Warn("предупреждение")                         // want `log message must be in English only`

	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck
	logger.Info("сервер запущен")                       // want `log message must be in English only`
	logger.Error("ошибка соединения")                   // want `log message must be in English only`
}
