package rules

import (
	"fmt"
	"go/ast"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var nonASCIILetterRanges = []*unicode.RangeTable{
	unicode.Cyrillic,
	unicode.Han,
	unicode.Hiragana,
	unicode.Katakana,
	unicode.Arabic,
	unicode.Hebrew,
	unicode.Devanagari,
	unicode.Greek,
	unicode.Thai,
	unicode.Hangul,
}

func CheckEnglishOnly(pass *analysis.Pass, msgArg ast.Expr, msg string) {
	for _, r := range msg {
		if isNonLatinLetter(r) {
			pass.Report(analysis.Diagnostic{
				Pos:      msgArg.Pos(),
				End:      msgArg.End(),
				Category: "auditlog-english",
				Message:  fmt.Sprintf("log message must be in English only, found non-Latin character %q", string(r)),
			})
			return
		}
	}
}

func IsValidEnglishOnly(msg string) bool {
	for _, r := range msg {
		if isNonLatinLetter(r) {
			return false
		}
	}
	return true
}

func isNonLatinLetter(r rune) bool {
	if !unicode.IsLetter(r) {
		return false
	}
	if r <= 0x7F {
		return false
	}
	return unicode.In(r, nonASCIILetterRanges...)
}
