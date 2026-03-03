# audit-log-linter

Статический анализатор для Go, совместимый с [golangci-lint](https://golangci-lint.run/). Проверяет лог-сообщения в вызовах `log/slog` и `go.uber.org/zap` на соответствие правилам оформления и безопасности.

## Правила

| ID | Правило | Пример нарушения |
|----|---------|-----------------|
| `auditlog-lowercase` | Сообщение должно начинаться со строчной буквы | `slog.Info("Starting server")` |
| `auditlog-english` | Сообщение должно быть только на английском | `slog.Error("ошибка подключения")` |
| `auditlog-special-chars` | Нельзя использовать спецсимволы и эмодзи | `slog.Info("started! 🚀")` |
| `auditlog-sensitive` | Нельзя выводить потенциально чувствительные данные | `slog.Info("token: " + token)` |

## Поддерживаемые логгеры

- [`log/slog`](https://pkg.go.dev/log/slog) — стандартная библиотека (Go 1.21+)
- [`go.uber.org/zap`](https://pkg.go.dev/go.uber.org/zap) — `Logger` и `SugaredLogger`

## Установка

### Standalone-бинарь

```bash
go install github.com/anvarulugov/audit-log-linter/cmd/auditloglint@latest
```

Запуск:

```bash
auditloglint ./...
```

### Плагин для golangci-lint (Module Plugin System)

1. Установить `golangci-lint` v2.x:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

2. Создать `.custom-gcl.yml` в корне проекта:

```yaml
version: v2.1.6
plugins:
  - module: "github.com/anvarulugov/audit-log-linter"
    import: "github.com/anvarulugov/audit-log-linter/plugin"
    version: v1.0.0
```

3. Создать или обновить `.golangci.yml`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - auditlog

linters-settings:
  custom:
    auditlog:
      type: "module"
      description: "Проверка лог-сообщений на соответствие правилам оформления."
      settings:
        disable_lowercase: false
        disable_english_only: false
        disable_no_special_chars: false
        disable_no_sensitive_data: false
        sensitive_keywords:
          - "internal_token"
```

4. Собрать кастомный бинарь и запустить:

```bash
golangci-lint custom      # создаёт ./custom-gcl
./custom-gcl run ./...
```

## Сборка из исходников

```bash
git clone https://github.com/anvarulugov/audit-log-linter
cd audit-log-linter
go build ./...
```

Сборка standalone-бинаря:

```bash
go build -o auditloglint ./cmd/auditloglint
```

## Тесты

```bash
go test ./...
```

## Конфигурация

Все настройки передаются через `.golangci.yml` в секции `linters-settings.custom.auditlog.settings`:

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|--------------|----------|
| `disable_lowercase` | bool | `false` | Отключить проверку строчной буквы |
| `disable_english_only` | bool | `false` | Отключить проверку английского языка |
| `disable_no_special_chars` | bool | `false` | Отключить проверку спецсимволов |
| `disable_no_sensitive_data` | bool | `false` | Отключить проверку чувствительных данных |
| `sensitive_keywords` | []string | `[]` | Дополнительные чувствительные ключевые слова |

### Встроенные чувствительные ключевые слова

`password`, `passwd`, `api_key`, `apikey`, `token`, `secret`, `credentials`, `auth`, `authorization`, `private_key`, `access_key`, `session`, `jwt`, `bearer`, `ssn`, `credit_card`, `cvv`, `pin`

## Примеры

### Корректные лог-сообщения

```go
slog.Info("starting server on port 8080")
slog.Error("failed to connect to database")
logger.Info("server started")
sugar.Infow("request completed", "duration", time.Since(start))
```

### Нарушения

```go
slog.Info("Starting server")              // auditlog-lowercase: заглавная буква
slog.Error("ошибка подключения")          // auditlog-english: не английский
slog.Info("server started!")              // auditlog-special-chars: символ '!'
slog.Warn("something went wrong...")      // auditlog-special-chars: многоточие
slog.Info("user password: " + password)  // auditlog-sensitive: конкатенация с sensitive-переменной
slog.Debug("api_key=" + apiKey)           // auditlog-sensitive: sensitive-префикс в конкатенации
```

## Структура проекта

```
.
├── analyzer/
│   ├── analyzer.go          # Основной analysis.Analyzer
│   ├── analyzer_test.go     # Интеграционные тесты (analysistest)
│   ├── config.go            # Структура конфигурации
│   ├── detector.go          # Определение вызовов log/slog и zap
│   └── rules/
│       ├── lowercase.go     # Правило 1: строчная буква
│       ├── english.go       # Правило 2: только английский
│       ├── special_chars.go # Правило 3: нет спецсимволов
│       └── sensitive.go     # Правило 4: нет чувствительных данных
├── cmd/auditloglint/
│   └── main.go              # Standalone CLI
├── plugin/
│   └── plugin.go            # Плагин для golangci-lint
├── testdata/
│   └── src/auditloglint/
│       ├── ok/              # Валидный код (диагностик не ожидается)
│       └── fail/            # Невалидный код (ожидаемые диагностики)
├── .custom-gcl.yml          # Конфигурация сборки кастомного golangci-lint
├── .golangci.example.yml    # Пример конфигурации golangci-lint с плагином
└── .github/workflows/ci.yml # GitHub Actions CI
```

## CI/CD

GitHub Actions запускает три джобы при каждом пуше и PR:

- `test` — `go test -race ./...`
- `lint` — `golangci-lint run`
- `build` — сборка standalone-бинаря и публикация артефакта
