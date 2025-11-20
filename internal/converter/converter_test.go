package converter

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidateInputFile tests the input file validation logic.
func TestValidateInputFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		wantError bool
		setup     func(string) error
	}{
		{
			name:      "empty path",
			filePath:  "",
			wantError: true,
		},
		{
			name:      "nonexistent file",
			filePath:  "/tmp/veve_nonexistent_" + randomString(10) + ".md",
			wantError: true,
		},
		{
			name:      "is a directory",
			filePath:  "/tmp",
			wantError: true,
		},
		{
			name:      "valid file",
			filePath:  filepath.Join(os.TempDir(), "veve_test_input.md"),
			wantError: false,
			setup: func(path string) error {
				return os.WriteFile(path, []byte("# Test"), 0o644)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(tt.filePath); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				defer os.Remove(tt.filePath)
			}

			err := ValidateInputFile(tt.filePath)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateInputFile() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestResolveOutputPath tests the output path resolution logic.
func TestResolveOutputPath(t *testing.T) {
	tests := []struct {
		inputPath  string
		outputPath string
		want       string
	}{
		{
			inputPath:  "/path/to/document.md",
			outputPath: "",
			want:       "/path/to/document.pdf",
		},
		{
			inputPath:  "README.markdown",
			outputPath: "",
			want:       "README.pdf",
		},
		{
			inputPath:  "noextension",
			outputPath: "",
			want:       "noextension.pdf",
		},
		{
			inputPath:  "/path/to/file.md",
			outputPath: "/custom/output.pdf",
			want:       "/custom/output.pdf",
		},
	}

	for _, tt := range tests {
		got := ResolveOutputPath(tt.inputPath, tt.outputPath)
		if got != tt.want {
			t.Errorf("ResolveOutputPath(%q, %q) = %q, want %q", tt.inputPath, tt.outputPath, got, tt.want)
		}
	}
}

// TestEnsureOutputDirectory tests the output directory creation logic.
func TestEnsureOutputDirectory(t *testing.T) {
	tests := []struct {
		name       string
		outputPath string
		wantError  bool
	}{
		{
			name:       "nested directories",
			outputPath: filepath.Join(os.TempDir(), "veve_test", "nested", "deep", "output.pdf"),
			wantError:  false,
		},
		{
			name:       "current directory",
			outputPath: "output.pdf",
			wantError:  false,
		},
		{
			name:       "dot directory",
			outputPath: "./output.pdf",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any previous test files
			defer os.RemoveAll(filepath.Dir(tt.outputPath))

			err := EnsureOutputDirectory(tt.outputPath)
			if (err != nil) != tt.wantError {
				t.Errorf("EnsureOutputDirectory() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				dir := filepath.Dir(tt.outputPath)
				if dir != "" && dir != "." {
					if _, err := os.Stat(dir); err != nil {
						t.Errorf("expected directory to be created, but got error: %v", err)
					}
				}
			}
		})
	}
}

// Helper to generate random strings for unique test file names
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
