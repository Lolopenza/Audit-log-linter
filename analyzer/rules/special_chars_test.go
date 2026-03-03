package rules_test

import (
	"testing"

	"github.com/anvarulugov/audit-log-linter/analyzer/rules"
)

func TestIsValidNoSpecialChars(t *testing.T) {
	cases := []struct {
		msg   string
		valid bool
	}{
		{"server started", true},
		{"connection failed", true},
		{"something went wrong", true},
		{"retrying in 5s", true},
		{"failed after 3 retries", true},
		{"", true},
		// Forbidden ASCII punctuation
		{"connection failed!!!", false},
		{"warning: something went wrong", false},
		{"server started!", false},
		// Ellipsis
		{"something went wrong...", false},
		{"loading...", false},
		// Emoji
		{"server started \U0001F680", false}, // 🚀
		{"connection failed \U0001F4A5", false}, // 💥
		// Allowed characters
		{"user id-123 processed", true},
		{"item, processed", true},
		{"item processed ok", true},
	}

	for _, tc := range cases {
		got := rules.IsValidNoSpecialChars(tc.msg)
		if got != tc.valid {
			t.Errorf("IsValidNoSpecialChars(%q) = %v, want %v", tc.msg, got, tc.valid)
		}
	}
}
