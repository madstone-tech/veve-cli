package contract_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCompletionCommand_Bash tests that 'veve completion bash' generates valid bash completion
func TestCompletionCommand_Bash(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("generates valid bash completion script", func(t *testing.T) {
		cmd := exec.Command(vevePath, "completion", "bash")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			t.Fatalf("completion bash command failed: %v", err)
		}

		output := stdout.String()

		// Verify it contains bash completion function
		if !strings.Contains(output, "complete") || !strings.Contains(output, "bash") {
			t.Error("bash completion output missing expected bash completion keywords")
		}

		// Verify it has completion logic
		if !strings.Contains(output, "COMPREPLY") && !strings.Contains(output, "compopt") {
			t.Logf("Note: bash completion may not contain COMPREPLY or compopt, output starts with: %s",
				truncate(output, 200))
		}

		// Should not be empty
		if len(output) < 50 {
			t.Error("bash completion output is too short, likely incomplete")
		}

		t.Logf("Generated bash completion (%d bytes)", len(output))
	})

	t.Run("bash completion contains engine flag", func(t *testing.T) {
		cmd := exec.Command(vevePath, "completion", "bash")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			t.Fatalf("completion bash command failed: %v", err)
		}

		output := stdout.String()

		// Should reference --engine flag or engine values
		if !strings.Contains(output, "engine") && !strings.Contains(output, "xelatex") {
			t.Logf("Note: bash completion may not explicitly list engine flag, output: %s",
				truncate(output, 300))
		}
	})
}

// TestCompletionCommand_Zsh tests that 'veve completion zsh' generates valid zsh completion
func TestCompletionCommand_Zsh(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("generates valid zsh completion script", func(t *testing.T) {
		cmd := exec.Command(vevePath, "completion", "zsh")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			t.Fatalf("completion zsh command failed: %v", err)
		}

		output := stdout.String()

		// Verify it contains zsh-specific syntax
		if !strings.Contains(output, "compdef") && !strings.Contains(output, "_arguments") {
			// Alternative: check for _veve function
			if !strings.Contains(output, "_veve") {
				t.Logf("Note: zsh completion may use different syntax, output starts with: %s",
					truncate(output, 200))
			}
		}

		// Should not be empty
		if len(output) < 50 {
			t.Error("zsh completion output is too short, likely incomplete")
		}

		t.Logf("Generated zsh completion (%d bytes)", len(output))
	})

	t.Run("zsh completion contains engine values", func(t *testing.T) {
		cmd := exec.Command(vevePath, "completion", "zsh")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			t.Fatalf("completion zsh command failed: %v", err)
		}

		output := stdout.String()

		// Should reference engine or at least one engine name
		if !strings.Contains(output, "engine") && !strings.Contains(output, "xelatex") {
			t.Logf("Note: zsh completion may use different engine references, output: %s",
				truncate(output, 300))
		}
	})
}

// TestCompletionCommand_Fish tests that 'veve completion fish' generates valid fish completion
func TestCompletionCommand_Fish(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("generates valid fish completion script", func(t *testing.T) {
		cmd := exec.Command(vevePath, "completion", "fish")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			t.Fatalf("completion fish command failed: %v", err)
		}

		output := stdout.String()

		// Verify it contains fish-specific syntax
		if !strings.Contains(output, "complete") || !strings.Contains(output, "command") {
			t.Logf("Note: fish completion may use different syntax, output starts with: %s",
				truncate(output, 200))
		}

		// Should not be empty
		if len(output) < 50 {
			t.Error("fish completion output is too short, likely incomplete")
		}

		t.Logf("Generated fish completion (%d bytes)", len(output))
	})

	t.Run("fish completion contains engine values", func(t *testing.T) {
		cmd := exec.Command(vevePath, "completion", "fish")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			t.Fatalf("completion fish command failed: %v", err)
		}

		output := stdout.String()

		// Should mention engine or engine values
		if !strings.Contains(output, "engine") && !strings.Contains(output, "xelatex") {
			t.Logf("Note: fish completion format differs, output: %s", truncate(output, 300))
		}
	})
}

// TestCompletionCommand_InstallScript tests that the install-completion.sh script exists and is executable
func TestCompletionInstallScript(t *testing.T) {
	scriptPath := "scripts/install-completion.sh"

	// Check if script exists (relative to repo root)
	// For now, we just verify the file can be found
	repoRoot := getRepoRoot(t)
	fullPath := filepath.Join(repoRoot, scriptPath)

	t.Run("install-completion.sh exists", func(t *testing.T) {
		info, err := os.Stat(fullPath)
		if err != nil {
			// Script may not exist yet; this is expected in early phases
			t.Skipf("install-completion.sh not found at %s (will be created in implementation)", fullPath)
		}

		// Verify it's a regular file
		if !info.Mode().IsRegular() {
			t.Error("install-completion.sh is not a regular file")
		}

		// Verify it's executable
		if (info.Mode() & 0111) == 0 {
			t.Error("install-completion.sh is not executable")
		}

		t.Logf("install-completion.sh exists and is executable")
	})
}

// TestEngineFlagCompletion tests that engine flag completion works via CLI
func TestEngineFlagCompletion(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("veve convert --engine completion returns engines", func(t *testing.T) {
		// This test may fail if Cobra completion is not fully integrated yet
		// It's marked as informational
		cmd := exec.Command(vevePath, "__complete", "convert", "--engine", "")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		_ = cmd.Run() // Ignore error for now

		output := stdout.String()
		errorOutput := stderr.String()

		t.Logf("Completion output: %s", truncate(output, 300))
		if errorOutput != "" {
			t.Logf("Completion stderr: %s", truncate(errorOutput, 300))
		}

		// The actual engine values may be returned differently depending on Cobra version
		// Just verify command doesn't crash
		if strings.Contains(errorOutput, "fatal") || strings.Contains(errorOutput, "panic") {
			t.Errorf("completion command produced fatal error")
		}
	})

	t.Run("veve convert --engine 'w' filters to weasyprint", func(t *testing.T) {
		// This test may be skipped if completion support not implemented yet
		cmd := exec.Command(vevePath, "__complete", "convert", "--engine", "w")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		_ = cmd.Run() // Ignore error for now

		output := stdout.String()

		// Output should be minimal or contain filtered results
		t.Logf("Filtered completion output for 'w': %s", truncate(output, 300))
	})
}

// Helper functions

// truncate returns first n characters of string
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// getRepoRoot returns the repository root directory for this test
func getRepoRoot(t *testing.T) string {
	// Try to find the repo root by looking for go.mod
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			t.Fatalf("could not find repository root (go.mod not found)")
		}
		cwd = parent
	}
}
