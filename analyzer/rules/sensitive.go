package rules

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var defaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"pass",
	"api_key",
	"apikey",
	"api-key",
	"token",
	"secret",
	"credentials",
	"credential",
	"auth",
	"authorization",
	"private_key",
	"privatekey",
	"private-key",
	"access_key",
	"accesskey",
	"access-key",
	"session",
	"jwt",
	"bearer",
	"ssn",
	"credit_card",
	"creditcard",
	"cvv",
	"pin",
}

func CheckNoSensitiveData(
	pass *analysis.Pass,
	msgArg ast.Expr,
	msg string,
	concatIdents []string,
	concatLiterals []string,
	extraKeywords []string,
) {
	if len(concatIdents) == 0 && len(concatLiterals) == 0 {
		return
	}

	keywords := mergeKeywords(defaultSensitiveKeywords, extraKeywords)

	for _, lit := range concatLiterals {
		if kw, found := containsSensitiveKeyword(lit, keywords); found {
			pass.Report(analysis.Diagnostic{
				Pos:      msgArg.Pos(),
				End:      msgArg.End(),
				Category: "auditlog-sensitive",
				Message:  fmt.Sprintf("log message may expose sensitive data: contains %q", kw),
			})
			return
		}
	}

	for _, ident := range concatIdents {
		if kw, found := containsSensitiveKeyword(ident, keywords); found {
			pass.Report(analysis.Diagnostic{
				Pos:      msgArg.Pos(),
				End:      msgArg.End(),
				Category: "auditlog-sensitive",
				Message:  fmt.Sprintf("log message may expose sensitive data via variable %q (matches keyword %q)", ident, kw),
			})
			return
		}
	}
}

func IsValidNoSensitiveData(concatIdents, concatLiterals, extraKeywords []string) bool {
	if len(concatIdents) == 0 && len(concatLiterals) == 0 {
		return true
	}
	keywords := mergeKeywords(defaultSensitiveKeywords, extraKeywords)
	for _, lit := range concatLiterals {
		if _, found := containsSensitiveKeyword(lit, keywords); found {
			return false
		}
	}
	for _, ident := range concatIdents {
		if _, found := containsSensitiveKeyword(ident, keywords); found {
			return false
		}
	}
	return true
}

func containsSensitiveKeyword(text string, keywords []string) (string, bool) {
	lower := strings.ToLower(text)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return kw, true
		}
	}
	return "", false
}

func mergeKeywords(defaults, extras []string) []string {
	seen := make(map[string]bool, len(defaults)+len(extras))
	result := make([]string, 0, len(defaults)+len(extras))
	for _, kw := range defaults {
		k := strings.ToLower(kw)
		if !seen[k] {
			seen[k] = true
			result = append(result, k)
		}
	}
	for _, kw := range extras {
		k := strings.ToLower(kw)
		if !seen[k] {
			seen[k] = true
			result = append(result, k)
		}
	}
	return result
}
