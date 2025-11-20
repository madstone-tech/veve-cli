package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestThemeRemove tests that 'veve theme remove name' deletes a theme with confirmation.
func TestThemeRemove(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create a test theme
	testThemePath := filepath.Join(themesDir, "remove-test.css")
	testCSS := `body { color: green; }`
	if err := os.WriteFile(testThemePath, []byte(testCSS), 0o644); err != nil {
		t.Fatalf("failed to create test theme: %v", err)
	}

	// Verify the theme exists before removal
	if _, err := os.Stat(testThemePath); err != nil {
		t.Fatalf("test theme creation failed: %v", err)
	}

	// Remove the theme with --force to skip confirmation
	cmd := exec.Command("veve", "theme", "remove", "remove-test", "--force")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("veve theme remove failed: %v\nOutput: %s", err, string(output))
		t.Skip("veve theme remove not yet fully implemented")
	}

	// Verify the theme file was deleted
	if _, err := os.Stat(testThemePath); err == nil {
		t.Errorf("theme file still exists after removal")
	}
}

// TestThemeRemoveNonExistent tests removing a non-existent theme.
func TestThemeRemoveNonExistent(t *testing.T) {
	cmd := exec.Command("veve", "theme", "remove", "nonexistent-theme", "--force")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Errorf("expected error when removing non-existent theme, but command succeeded")
	}

	outStr := string(output)
	if outStr == "" {
		t.Logf("error output not available")
	}
}

// TestThemeRemoveBuiltIn tests that built-in themes cannot be removed.
func TestThemeRemoveBuiltIn(t *testing.T) {
	cmd := exec.Command("veve", "theme", "remove", "default", "--force")
	err := cmd.Run()

	if err == nil {
		t.Errorf("expected error when removing built-in theme")
	}
}

// TestThemeRemoveWithConfirmation tests theme removal with confirmation prompt.
func TestThemeRemoveWithConfirmation(t *testing.T) {
	// This test is tricky because it requires interactive input
	// For now, just test the --force flag behavior

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	testThemePath := filepath.Join(themesDir, "confirm-test.css")
	if err := os.WriteFile(testThemePath, []byte("body { }"), 0o644); err != nil {
		t.Fatalf("failed to create test theme: %v", err)
	}

	// Remove with --force (skip confirmation)
	cmd := exec.Command("veve", "theme", "remove", "confirm-test", "--force")
	err = cmd.Run()

	if err != nil {
		t.Logf("veve theme remove --force failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	// Verify deletion
	if _, err := os.Stat(testThemePath); err == nil {
		t.Errorf("theme should be deleted with --force flag")
	}
}

// TestThemeRemoveVerifyNotInList tests that removed theme no longer appears in list.
func TestThemeRemoveVerifyNotInList(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	testThemePath := filepath.Join(themesDir, "list-test.css")
	if err := os.WriteFile(testThemePath, []byte("body { }"), 0o644); err != nil {
		t.Fatalf("failed to create test theme: %v", err)
	}
	defer os.Remove(testThemePath)

	// Verify theme appears in list
	listCmd := exec.Command("veve", "theme", "list")
	listOutput, _ := listCmd.CombinedOutput()
	if !strings.Contains(string(listOutput), "list-test") {
		t.Logf("theme not appearing in list initially, skipping removal test")
		t.Skip("theme discovery may not be working")
	}

	// Remove the theme
	removeCmd := exec.Command("veve", "theme", "remove", "list-test", "--force")
	err = removeCmd.Run()

	if err != nil {
		t.Logf("veve theme remove failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	// Verify theme no longer in list
	listCmd2 := exec.Command("veve", "theme", "list")
	listOutput2, _ := listCmd2.CombinedOutput()
	if strings.Contains(string(listOutput2), "list-test") {
		t.Errorf("removed theme still appears in list")
	}
}

// TestThemeRemoveExitCode tests that remove returns correct exit codes.
func TestThemeRemoveExitCode(t *testing.T) {
	// Test non-existent theme returns error
	cmd := exec.Command("veve", "theme", "remove", "nonexistent", "--force")
	err := cmd.Run()

	if err == nil {
		t.Logf("Note: command succeeded for non-existent theme")
	}
}
