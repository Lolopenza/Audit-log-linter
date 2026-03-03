package main

import (
	"flag"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/anvarulugov/audit-log-linter/analyzer"
)

func main() {
	var cfg analyzer.Config

	flag.BoolVar(&cfg.DisableLowercase, "disable_lowercase", false, "disable check: message must start with lowercase")
	flag.BoolVar(&cfg.DisableEnglishOnly, "disable_english_only", false, "disable check: message must be English only")
	flag.BoolVar(&cfg.DisableNoSpecialChars, "disable_no_special_chars", false, "disable check: no special chars/emoji")
	flag.BoolVar(&cfg.DisableNoSensitiveData, "disable_no_sensitive_data", false, "disable check: no sensitive data")

	singlechecker.Main(analyzer.NewAnalyzer(cfg))
}
