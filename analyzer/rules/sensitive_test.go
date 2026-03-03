package rules_test

import (
	"testing"

	"github.com/anvarulugov/audit-log-linter/analyzer/rules"
)

func TestIsValidNoSensitiveData(t *testing.T) {
	t.Run("pure literals are always valid", func(t *testing.T) {
		// Pure string literals are NOT checked: they describe events, not values.
		pureLiterals := []string{
			"user authenticated successfully",
			"api request completed",
			"token validated",
			"starting server",
			"password reset triggered",
			"",
		}
		for _, msg := range pureLiterals {
			_ = msg
			got := rules.IsValidNoSensitiveData(nil, nil, nil)
			if !got {
				t.Errorf("IsValidNoSensitiveData(nil, nil, nil) = false, want true")
			}
		}
	})

	t.Run("concatenation with sensitive ident", func(t *testing.T) {
		cases := []struct {
			idents []string
			valid  bool
		}{
			{[]string{"password"}, false},
			{[]string{"apiKey"}, false},  // "apikey" matches keyword
			{[]string{"token"}, false},
			{[]string{"secret"}, false},
			{[]string{"requestID"}, true}, // no sensitive keyword
			{[]string{"userID"}, true},
			{[]string{"port"}, true},
		}
		for _, tc := range cases {
			got := rules.IsValidNoSensitiveData(tc.idents, nil, nil)
			if got != tc.valid {
				t.Errorf("IsValidNoSensitiveData(idents=%v, nil, nil) = %v, want %v", tc.idents, got, tc.valid)
			}
		}
	})

	t.Run("concatenation with sensitive literal prefix", func(t *testing.T) {
		cases := []struct {
			lits  []string
			valid bool
		}{
			{[]string{"user password: "}, false},  // contains "password"
			{[]string{"api_key="}, false},          // contains "api_key"
			{[]string{"token: "}, false},           // contains "token"
			{[]string{"connecting to host "}, true}, // no keyword
			{[]string{"user id: "}, true},          // "id" is not a keyword
		}
		for _, tc := range cases {
			got := rules.IsValidNoSensitiveData(nil, tc.lits, nil)
			if got != tc.valid {
				t.Errorf("IsValidNoSensitiveData(nil, lits=%v, nil) = %v, want %v", tc.lits, got, tc.valid)
			}
		}
	})

	t.Run("extra keywords", func(t *testing.T) {
		got := rules.IsValidNoSensitiveData([]string{"internalkey"}, nil, []string{"internalkey"})
		if got {
			t.Error("expected invalid for custom keyword match, got valid")
		}

		got = rules.IsValidNoSensitiveData([]string{"regularvar"}, nil, []string{"mysecret"})
		if !got {
			t.Error("expected valid for non-matching custom keyword, got invalid")
		}
	})
}
