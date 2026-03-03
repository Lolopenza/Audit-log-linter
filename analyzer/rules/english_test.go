package rules_test

import (
	"testing"

	"github.com/anvarulugov/audit-log-linter/analyzer/rules"
)

func TestIsValidEnglishOnly(t *testing.T) {
	cases := []struct {
		msg   string
		valid bool
	}{
		{"starting server on port 8080", true},
		{"failed to connect to database", true},
		{"user authenticated successfully", true},
		{"api request completed", true},
		{"", true},
		{"test 123 ok", true},
		{"запуск сервера", false},           // Cyrillic
		{"ошибка подключения", false},       // Cyrillic
		{"サーバー起動", false},              // Japanese
		{"连接失败", false},                 // Chinese
		{"mixed english и русский", false}, // Mixed
	}

	for _, tc := range cases {
		got := rules.IsValidEnglishOnly(tc.msg)
		if got != tc.valid {
			t.Errorf("IsValidEnglishOnly(%q) = %v, want %v", tc.msg, got, tc.valid)
		}
	}
}
