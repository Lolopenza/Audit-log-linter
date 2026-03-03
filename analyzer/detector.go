package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/types/typeutil"
)

var supportedLogPkgs = map[string]bool{
	"log/slog":        true,
	"go.uber.org/zap": true,
}

var slogMsgArgIndex = map[string]int{
	"Debug":        0,
	"Info":         0,
	"Warn":         0,
	"Error":        0,
	"Log":          2,
	"DebugContext": 1,
	"InfoContext":  1,
	"WarnContext":  1,
	"ErrorContext": 1,
}

var zapLogFunctions = map[string]bool{
	"Debug":   true,
	"Info":    true,
	"Warn":    true,
	"Error":   true,
	"DPanic":  true,
	"Panic":   true,
	"Fatal":   true,
	"Debugw":  true,
	"Infow":   true,
	"Warnw":   true,
	"Errorw":  true,
	"DPanicw": true,
	"Panicw":  true,
	"Fatalw":  true,
	"Debugf":  true,
	"Infof":   true,
	"Warnf":   true,
	"Errorf":  true,
	"DPanicf": true,
	"Panicf":  true,
	"Fatalf":  true,
}

type LogCall struct {
	Call        *ast.CallExpr
	MsgArg      ast.Expr
	MsgLit      string
	MsgArgIndex int
}

func DetectLogCall(typesInfo *types.Info, call *ast.CallExpr) (LogCall, bool) {
	if typesInfo == nil {
		return LogCall{}, false
	}

	fn := typeutil.StaticCallee(typesInfo, call)
	if fn == nil {
		return LogCall{}, false
	}

	pkg := fn.Pkg()
	if pkg == nil {
		return LogCall{}, false
	}

	pkgPath := pkg.Path()
	if !supportedLogPkgs[pkgPath] {
		return LogCall{}, false
	}

	fnName := fn.Name()
	msgIdx := -1

	switch pkgPath {
	case "log/slog":
		idx, ok := slogMsgArgIndex[fnName]
		if !ok {
			return LogCall{}, false
		}
		msgIdx = idx
	case "go.uber.org/zap":
		if !zapLogFunctions[fnName] {
			return LogCall{}, false
		}
		msgIdx = 0
	}

	if msgIdx < 0 || len(call.Args) <= msgIdx {
		return LogCall{}, false
	}

	msgArg := call.Args[msgIdx]
	msgLit := extractStringLiteral(typesInfo, msgArg)

	return LogCall{
		Call:        call,
		MsgArg:      msgArg,
		MsgLit:      msgLit,
		MsgArgIndex: msgIdx,
	}, true
}

func extractStringLiteral(typesInfo *types.Info, expr ast.Expr) string {
	if typesInfo == nil {
		return ""
	}
	tv, ok := typesInfo.Types[expr]
	if !ok {
		return ""
	}
	if tv.Value == nil || tv.Value.Kind() != constant.String {
		return ""
	}
	return constant.StringVal(tv.Value)
}

type SensitivePart struct {
	Literal         string
	Ident           string
	IsConcatenation bool
}

func ExtractSensitiveParts(typesInfo *types.Info, expr ast.Expr) []SensitivePart {
	var parts []SensitivePart
	collectParts(typesInfo, expr, &parts)
	return parts
}

func collectParts(typesInfo *types.Info, expr ast.Expr, parts *[]SensitivePart) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			collectParts(typesInfo, e.X, parts)
			collectParts(typesInfo, e.Y, parts)
		}
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			val := strings.Trim(e.Value, `"`)
			val = strings.ReplaceAll(val, "`", "")
			*parts = append(*parts, SensitivePart{Literal: val, IsConcatenation: true})
		}
	case *ast.Ident:
		*parts = append(*parts, SensitivePart{Ident: e.Name, IsConcatenation: true})
	case *ast.SelectorExpr:
		*parts = append(*parts, SensitivePart{Ident: e.Sel.Name, IsConcatenation: true})
	default:
		if lit := extractStringLiteral(typesInfo, expr); lit != "" {
			*parts = append(*parts, SensitivePart{Literal: lit, IsConcatenation: true})
		}
	}
}

func IsConcatenationExpr(expr ast.Expr) bool {
	b, ok := expr.(*ast.BinaryExpr)
	return ok && b.Op == token.ADD
}
