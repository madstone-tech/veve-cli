// Package engines provides PDF engine detection, validation, and selection logic.
package engines

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DetectInstalledEngines searches PATH for available PDF engines
// Returns a slice of PDFEngine with IsInstalled set based on availability
func DetectInstalledEngines() ([]PDFEngine, error) {
	definitions := DefaultEngineDefinitions()
	var installed []PDFEngine

	for _, name := range PriorityOrder {
		def, exists := definitions[name]
		if !exists {
			continue
		}

		// Check if engine binary is in PATH
		if _, err := exec.LookPath(name); err == nil {
			def.IsInstalled = true

			// Try to detect version
			if version, err := getEngineVersion(name); err == nil {
				def.Version = version
			}

			installed = append(installed, def)
		}
	}

	// If no engines found, return error with helpful message
	if len(installed) == 0 {
		return nil, fmt.Errorf("no PDF rendering engines found in PATH; " +
			"please install one of: xelatex, lualatex, weasyprint, or prince")
	}

	return installed, nil
}

// getEngineVersion attempts to detect engine version
// Returns version string or error if detection fails
func getEngineVersion(engineName string) (string, error) {
	switch engineName {
	case "xelatex", "lualatex":
		// LaTeX engines: try --version
		cmd := exec.Command(engineName, "--version")
		output, err := cmd.CombinedOutput()
		if err == nil {
			// Extract version from first line
			lines := strings.Split(string(output), "\n")
			if len(lines) > 0 && lines[0] != "" {
				return lines[0], nil
			}
		}

	case "weasyprint":
		// WeasyPrint: try --version
		cmd := exec.Command("weasyprint", "--version")
		output, err := cmd.CombinedOutput()
		if err == nil {
			version := strings.TrimSpace(string(output))
			if version != "" {
				return version, nil
			}
		}

	case "prince":
		// Prince: try --version
		cmd := exec.Command("prince", "--version")
		output, err := cmd.CombinedOutput()
		if err == nil {
			version := strings.TrimSpace(string(output))
			if version != "" {
				return version, nil
			}
		}
	}

	return "", fmt.Errorf("could not detect version for %s", engineName)
}

// FindEngineInPath searches for a specific engine binary in system PATH
// Returns the full path to the engine or error if not found
func FindEngineInPath(engineName string) (string, error) {
	path, err := exec.LookPath(engineName)
	if err != nil {
		return "", fmt.Errorf("engine '%s' not found in PATH", engineName)
	}

	// Verify it's actually an executable file
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("cannot stat engine at %s: %w", path, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("engine path is a directory: %s", path)
	}

	return path, nil
}

// GetAllPaths returns the system PATH as a slice of directories
func GetAllPaths() []string {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil
	}

	return filepath.SplitList(pathEnv)
}

// SearchEngineInDirs searches for engine in specific directories
// Returns full path or error if not found
func SearchEngineInDirs(engineName string, dirs []string) (string, error) {
	for _, dir := range dirs {
		enginePath := filepath.Join(dir, engineName)

		// On Windows, try with .exe extension
		if _, err := os.Stat(enginePath); err == nil {
			return enginePath, nil
		}

		// Try with .exe (for Windows or when running on Windows)
		exePath := enginePath + ".exe"
		if _, err := os.Stat(exePath); err == nil {
			return exePath, nil
		}
	}

	return "", fmt.Errorf("engine '%s' not found in specified directories", engineName)
}
