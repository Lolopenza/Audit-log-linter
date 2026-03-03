package rules

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func CheckLowercase(pass *analysis.Pass, msgArg ast.Expr, msg string) {
	if msg == "" {
		return
	}

	firstRune := []rune(msg)[0]
	if !unicode.IsUpper(firstRune) {
		return
	}

	lowered := strings.ToLower(string(firstRune)) + msg[len(string(firstRune)):]
	pos := msgArg.Pos()
	end := msgArg.End()

	pass.Report(analysis.Diagnostic{
		Pos:      pos,
		End:      end,
		Category: "auditlog-lowercase",
		Message:  fmt.Sprintf("log message must start with a lowercase letter, got %q", string(firstRune)),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("change %q to %q", string(firstRune), strings.ToLower(string(firstRune))),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     pos,
						End:     end,
						NewText: []byte(fmt.Sprintf("%q", lowered)),
					},
				},
			},
		},
	})
}

func IsValidLowercase(msg string) bool {
	if msg == "" {
		return true
	}
	firstRune := []rune(msg)[0]
	return !unicode.IsUpper(firstRune)
}
