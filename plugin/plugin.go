package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/anvarulugov/audit-log-linter/analyzer"
)

func init() {
	register.Plugin("auditlog", New)
}

type PluginSettings struct {
	DisableLowercase       bool     `json:"disable_lowercase"         mapstructure:"disable_lowercase"`
	DisableEnglishOnly     bool     `json:"disable_english_only"      mapstructure:"disable_english_only"`
	DisableNoSpecialChars  bool     `json:"disable_no_special_chars"  mapstructure:"disable_no_special_chars"`
	DisableNoSensitiveData bool     `json:"disable_no_sensitive_data" mapstructure:"disable_no_sensitive_data"`
	SensitiveKeywords      []string `json:"sensitive_keywords"        mapstructure:"sensitive_keywords"`
}

type auditlogPlugin struct {
	cfg analyzer.Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[PluginSettings](settings)
	if err != nil {
		return nil, err
	}

	cfg := analyzer.Config{
		DisableLowercase:       s.DisableLowercase,
		DisableEnglishOnly:     s.DisableEnglishOnly,
		DisableNoSpecialChars:  s.DisableNoSpecialChars,
		DisableNoSensitiveData: s.DisableNoSensitiveData,
		SensitiveKeywords:      s.SensitiveKeywords,
	}

	return &auditlogPlugin{cfg: cfg}, nil
}

func (p *auditlogPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.NewAnalyzer(p.cfg),
	}, nil
}

func (p *auditlogPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
