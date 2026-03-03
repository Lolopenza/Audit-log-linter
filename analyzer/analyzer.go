package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/anvarulugov/audit-log-linter/analyzer/rules"
)

const analyzerName = "auditlog"
const analyzerDoc = `checks log messages in log/slog and go.uber.org/zap calls for style violations

Rules:
  1. auditlog-lowercase:      messages must start with a lowercase letter
  2. auditlog-english:        messages must be in English (no non-Latin Unicode letters)
  3. auditlog-special-chars:  messages must not contain special characters or emoji
  4. auditlog-sensitive:      messages must not expose sensitive data (password, token, etc.)`

func NewAnalyzer(cfg Config) *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name:     analyzerName,
		Doc:      analyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	a.Run = func(pass *analysis.Pass) (any, error) {
		return run(pass, cfg)
	}
	return a
}

var Analyzer = NewAnalyzer(DefaultConfig())

func run(pass *analysis.Pass, cfg Config) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(node ast.Node) {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}

		logCall, found := DetectLogCall(pass.TypesInfo, call)
		if !found {
			return
		}

		msg := logCall.MsgLit
		msgArg := logCall.MsgArg

		if !cfg.DisableLowercase && msg != "" {
			rules.CheckLowercase(pass, msgArg, msg)
		}

		if !cfg.DisableEnglishOnly && msg != "" {
			rules.CheckEnglishOnly(pass, msgArg, msg)
		}

		if !cfg.DisableNoSpecialChars && msg != "" {
			rules.CheckNoSpecialChars(pass, msgArg, msg)
		}

		if !cfg.DisableNoSensitiveData {
			var concatIdents, concatLiterals []string
			if IsConcatenationExpr(msgArg) {
				for _, part := range ExtractSensitiveParts(pass.TypesInfo, msgArg) {
					if part.Ident != "" {
						concatIdents = append(concatIdents, part.Ident)
					}
					if part.Literal != "" {
						concatLiterals = append(concatLiterals, part.Literal)
					}
				}
			}
			rules.CheckNoSensitiveData(pass, msgArg, msg, concatIdents, concatLiterals, cfg.SensitiveKeywords)
		}
	})

	return nil, nil
}
