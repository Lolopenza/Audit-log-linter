package analyzer

type Config struct {
	DisableLowercase       bool     `json:"disable_lowercase"         mapstructure:"disable_lowercase"`
	DisableEnglishOnly     bool     `json:"disable_english_only"      mapstructure:"disable_english_only"`
	DisableNoSpecialChars  bool     `json:"disable_no_special_chars"  mapstructure:"disable_no_special_chars"`
	DisableNoSensitiveData bool     `json:"disable_no_sensitive_data" mapstructure:"disable_no_sensitive_data"`
	SensitiveKeywords      []string `json:"sensitive_keywords"        mapstructure:"sensitive_keywords"`
}

func DefaultConfig() Config {
	return Config{}
}
