package rules_test

import (
	"testing"

	"github.com/anvarulugov/audit-log-linter/analyzer/rules"
)

func TestIsValidLowercase(t *testing.T) {
	cases := []struct {
		msg   string
		valid bool
	}{
		{"starting server", true},
		{"failed to connect", true},
		{"", true},          // empty string is valid
		{"Starting server", false},
		{"Failed to connect", false},
		{"ERROR: something", false},
		{"запуск", true},    // non-ASCII: not a Latin uppercase letter, skip
		{"a", true},
		{"A", false},
	}

	for _, tc := range cases {
		got := rules.IsValidLowercase(tc.msg)
		if got != tc.valid {
			t.Errorf("IsValidLowercase(%q) = %v, want %v", tc.msg, got, tc.valid)
		}
	}
}
