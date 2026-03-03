package analyzer_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/anvarulugov/audit-log-linter/analyzer"
)

func testdataDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file path")
	}
	// thisFile is .../analyzer/analyzer_test.go.
	// analysistest.Run expects the GOPATH root so that packages are resolved as
	// $dir/src/<pkg>. We set dir to the "testdata" directory (one level up from
	// the analyzer package, then into testdata/).
	projectRoot := filepath.Dir(filepath.Dir(thisFile))
	return filepath.Join(projectRoot, "testdata")
}

func TestAnalyzer_OK(t *testing.T) {
	analysistest.Run(t, testdataDir(t), analyzer.Analyzer, "auditloglint/ok")
}

func TestAnalyzer_Lowercase(t *testing.T) {
	analysistest.Run(t, testdataDir(t), analyzer.Analyzer, "auditloglint/fail/lowercase")
}

func TestAnalyzer_English(t *testing.T) {
	analysistest.Run(t, testdataDir(t), analyzer.Analyzer, "auditloglint/fail/english")
}

func TestAnalyzer_SpecialChars(t *testing.T) {
	analysistest.Run(t, testdataDir(t), analyzer.Analyzer, "auditloglint/fail/special_chars")
}

func TestAnalyzer_Sensitive(t *testing.T) {
	analysistest.Run(t, testdataDir(t), analyzer.Analyzer, "auditloglint/fail/sensitive")
}

func TestAnalyzer_DisabledRules(t *testing.T) {
	cfg := analyzer.Config{
		DisableLowercase:       true,
		DisableEnglishOnly:     true,
		DisableNoSpecialChars:  true,
		DisableNoSensitiveData: true,
	}
	a := analyzer.NewAnalyzer(cfg)
	// With all rules disabled, the "ok" package should still pass.
	analysistest.Run(t, testdataDir(t), a, "auditloglint/ok")
}

func TestAnalyzer_CustomSensitiveKeywords(t *testing.T) {
	cfg := analyzer.Config{
		SensitiveKeywords: []string{"mysecret", "internalkey"},
	}
	a := analyzer.NewAnalyzer(cfg)
	// Custom keywords are additive; the ok package should still pass.
	analysistest.Run(t, testdataDir(t), a, "auditloglint/ok")
}
