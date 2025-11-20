package config

import (
	"os"
	"path/filepath"
)

// Paths represents XDG Base Directory paths for veve configuration and data.
type Paths struct {
	// ConfigDir is the user's config directory (~/.config/veve on Unix, %APPDATA%/veve on Windows)
	ConfigDir string
	// DataDir is the user's data directory (~/.local/share/veve on Unix, %APPDATA%/veve on Windows)
	DataDir string
	// CacheDir is the user's cache directory (~/.cache/veve on Unix, %TEMP%/veve on Windows)
	CacheDir string
	// ThemesDir is the directory containing user themes
	ThemesDir string
	// ConfigFile is the main veve.toml configuration file path
	ConfigFile string
}

// GetPaths returns XDG Base Directory paths for the current platform.
// On Unix systems, it respects XDG_CONFIG_HOME, XDG_DATA_HOME, and XDG_CACHE_HOME environment variables.
// On Windows, it uses %APPDATA%.
func GetPaths() (Paths, error) {
	var (
		configDir string
		dataDir   string
		cacheDir  string
	)

	// Determine config directory
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		configDir = filepath.Join(xdgConfig, "veve")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return Paths{}, err
		}
		configDir = filepath.Join(home, ".config", "veve")
	}

	// Determine data directory
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		dataDir = filepath.Join(xdgData, "veve")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return Paths{}, err
		}
		dataDir = filepath.Join(home, ".local", "share", "veve")
	}

	// Determine cache directory
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		cacheDir = filepath.Join(xdgCache, "veve")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return Paths{}, err
		}
		cacheDir = filepath.Join(home, ".cache", "veve")
	}

	themesDir := filepath.Join(configDir, "themes")
	configFile := filepath.Join(configDir, "veve.toml")

	return Paths{
		ConfigDir:  configDir,
		DataDir:    dataDir,
		CacheDir:   cacheDir,
		ThemesDir:  themesDir,
		ConfigFile: configFile,
	}, nil
}

// EnsureDirectories creates all necessary veve directories if they don't exist.
func (p *Paths) EnsureDirectories() error {
	dirs := []string{p.ConfigDir, p.DataDir, p.CacheDir, p.ThemesDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}
