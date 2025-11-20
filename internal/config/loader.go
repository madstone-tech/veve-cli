package config

import (
	"github.com/spf13/viper"
)

// Config represents veve's configuration loaded from veve.toml.
type Config struct {
	// PDFEngine is the Pandoc PDF engine to use (default: "pdflatex")
	PDFEngine string `mapstructure:"pdf_engine"`
	// DefaultTheme is the default theme to use for conversions
	DefaultTheme string `mapstructure:"default_theme"`
	// Verbose enables verbose output
	Verbose bool `mapstructure:"verbose"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		PDFEngine:    "pdflatex",
		DefaultTheme: "default",
		Verbose:      false,
	}
}

// LoadConfig loads the veve configuration from veve.toml.
// If the config file doesn't exist, returns the default configuration.
func LoadConfig(configFile string) (Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Configure Viper
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("toml")

	// Set defaults
	v.SetDefault("pdf_engine", cfg.PDFEngine)
	v.SetDefault("default_theme", cfg.DefaultTheme)
	v.SetDefault("verbose", cfg.Verbose)

	// Try to read the config file (it's okay if it doesn't exist)
	if err := v.ReadInConfig(); err != nil {
		// It's fine if the file doesn't exist; we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Real error occurred
			return cfg, err
		}
	}

	// Unmarshal into our Config struct
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// SaveConfig saves the configuration to veve.toml.
func SaveConfig(configFile string, cfg Config) error {
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("toml")

	v.Set("pdf_engine", cfg.PDFEngine)
	v.Set("default_theme", cfg.DefaultTheme)
	v.Set("verbose", cfg.Verbose)

	return v.WriteConfigAs(configFile)
}
