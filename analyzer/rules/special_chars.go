package rules

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var forbiddenASCII = map[rune]bool{
	'!':  true,
	'@':  true,
	'#':  true,
	'$':  true,
	'^':  true,
	'&':  true,
	'*':  true,
	'+':  true,
	'=':  true,
	'{':  true,
	'}':  true,
	'[':  true,
	']':  true,
	'|':  true,
	'\\': true,
	';':  true,
	':':  true,
	'"':  true,
	'<':  true,
	'>':  true,
	'?':  true,
	'`':  true,
	'~':  true,
}

var emojiRanges = []*unicode.RangeTable{
	{R32: []unicode.Range32{{Lo: 0x1F600, Hi: 0x1F64F, Stride: 1}}},
	{R32: []unicode.Range32{{Lo: 0x1F300, Hi: 0x1F5FF, Stride: 1}}},
	{R32: []unicode.Range32{{Lo: 0x1F680, Hi: 0x1F6FF, Stride: 1}}},
	{R32: []unicode.Range32{{Lo: 0x1F900, Hi: 0x1F9FF, Stride: 1}}},
	{R32: []unicode.Range32{{Lo: 0x1FA00, Hi: 0x1FA6F, Stride: 1}}},
	{R32: []unicode.Range32{{Lo: 0x1FA70, Hi: 0x1FAFF, Stride: 1}}},
	{R16: []unicode.Range16{{Lo: 0x2702, Hi: 0x27B0, Stride: 1}}},
	{R16: []unicode.Range16{{Lo: 0x2600, Hi: 0x26FF, Stride: 1}}},
	{R32: []unicode.Range32{{Lo: 0x1F100, Hi: 0x1F1FF, Stride: 1}}},
	{R32: []unicode.Range32{{Lo: 0x1F1E0, Hi: 0x1F1FF, Stride: 1}}},
}

func CheckNoSpecialChars(pass *analysis.Pass, msgArg ast.Expr, msg string) {
	if strings.Contains(msg, "...") {
		pass.Report(analysis.Diagnostic{
			Pos:      msgArg.Pos(),
			End:      msgArg.End(),
			Category: "auditlog-special-chars",
			Message:  `log message must not contain "..." (ellipsis)`,
		})
		return
	}

	for _, r := range msg {
		if unicode.In(r, emojiRanges...) {
			pass.Report(analysis.Diagnostic{
				Pos:      msgArg.Pos(),
				End:      msgArg.End(),
				Category: "auditlog-special-chars",
				Message:  fmt.Sprintf("log message must not contain emoji: %q", string(r)),
			})
			return
		}

		if r <= 0x7F && forbiddenASCII[r] {
			pass.Report(analysis.Diagnostic{
				Pos:      msgArg.Pos(),
				End:      msgArg.End(),
				Category: "auditlog-special-chars",
				Message:  fmt.Sprintf("log message must not contain special character %q", string(r)),
			})
			return
		}
	}
}

func IsValidNoSpecialChars(msg string) bool {
	if strings.Contains(msg, "...") {
		return false
	}
	for _, r := range msg {
		if unicode.In(r, emojiRanges...) {
			return false
		}
		if r <= 0x7F && forbiddenASCII[r] {
			return false
		}
	}
	return true
}
