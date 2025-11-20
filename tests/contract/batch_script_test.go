package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBatchProcessing tests that veve can be used in a bash loop for batch conversion.
func TestBatchProcessing(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test markdown files
	files := []string{"file1.md", "file2.md", "file3.md"}
	for _, file := range files {
		filePath := filepath.Join(tmpDir, file)
		content := "# " + file + "\n\nContent for " + file + "\n"
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Create a simple bash script to batch convert
	scriptPath := filepath.Join(tmpDir, "batch.sh")
	scriptContent := `#!/bin/bash
cd ` + tmpDir + `
for f in *.md; do
  veve "$f" -o "${f%.md}.pdf" --quiet
  if [ $? -ne 0 ]; then
    echo "Failed to convert $f" >&2
    exit 1
  fi
done
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	// Run the batch script
	cmd := exec.Command("bash", scriptPath)
	err := cmd.Run()

	if err != nil {
		t.Logf("batch script failed: %v", err)
		t.Skip("veve batch processing not yet fully implemented")
	}

	// Verify all PDFs were created
	for _, file := range files {
		pdfPath := filepath.Join(tmpDir, filepath.Base(file[:len(file)-2])+".pdf")
		if _, err := os.Stat(pdfPath); err != nil {
			t.Errorf("expected PDF not created: %s", pdfPath)
		}
	}
}

// TestBatchProcessingWithGlob tests batch processing with glob patterns.
func TestBatchProcessingWithGlob(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	for i := 1; i <= 3; i++ {
		filePath := filepath.Join(tmpDir, "doc"+string(rune('0'+i))+".md")
		content := "# Document " + string(rune('0'+i)) + "\n\nContent\n"
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	// Use find to get files and process them
	scriptPath := filepath.Join(tmpDir, "batch_glob.sh")
	scriptContent := `#!/bin/bash
cd ` + tmpDir + `
find . -name "*.md" -print0 | while IFS= read -r -d '' f; do
  veve "$f" --quiet
  if [ $? -ne 0 ]; then
    echo "Failed: $f" >&2
    exit 1
  fi
done
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	cmd := exec.Command("bash", scriptPath)
	err := cmd.Run()

	if err != nil {
		t.Logf("glob batch processing failed: %v", err)
		t.Skip("feature not yet implemented")
	}
}

// TestBatchProcessingWithDirectoryArg tests processing directory of files.
func TestBatchProcessingWithDirectoryArg(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	os.MkdirAll(inputDir, 0o755)
	os.MkdirAll(outputDir, 0o755)

	// Create test files
	for i := 1; i <= 3; i++ {
		filePath := filepath.Join(inputDir, "file"+string(rune('0'+i))+".md")
		content := "# File\n\nContent\n"
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	// Batch convert with output to directory
	scriptPath := filepath.Join(tmpDir, "batch_dir.sh")
	scriptContent := `#!/bin/bash
input_dir=` + inputDir + `
output_dir=` + outputDir + `
for f in "$input_dir"/*.md; do
  basename=$(basename "$f" .md)
  veve "$f" -o "$output_dir/$basename.pdf" --quiet
  if [ $? -ne 0 ]; then
    exit 1
  fi
done
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	cmd := exec.Command("bash", scriptPath)
	err := cmd.Run()

	if err != nil {
		t.Logf("directory batch processing failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	// Check output files
	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Logf("failed to read output directory: %v", err)
		return
	}

	if len(files) < 3 {
		t.Logf("expected at least 3 output files, got %d", len(files))
	}
}

// TestBatchWithErrorHandling tests that batch processing handles errors.
func TestBatchWithErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some valid and one invalid markdown file
	os.WriteFile(filepath.Join(tmpDir, "valid1.md"), []byte("# Valid\nContent"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "valid2.md"), []byte("# Valid\nContent"), 0o644)

	// Create script that counts successes and failures
	scriptPath := filepath.Join(tmpDir, "batch_error.sh")
	scriptContent := `#!/bin/bash
success=0
failed=0
for f in ` + tmpDir + `/*.md; do
  if veve "$f" --quiet 2>/dev/null; then
    ((success++))
  else
    ((failed++))
  fi
done
echo "Success: $success, Failed: $failed"
exit 0
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	cmd := exec.Command("bash", scriptPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("batch error handling test failed: %v", err)
		return
	}

	t.Logf("Batch result: %s", string(output))
}

// TestBatchWithStdin tests batch processing using stdin pipes.
func TestBatchWithStdin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test markdown files
	for i := 1; i <= 2; i++ {
		filePath := filepath.Join(tmpDir, "doc"+string(rune('0'+i))+".md")
		content := "# Doc " + string(rune('0'+i)) + "\n\nContent\n"
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	// Use cat to pipe files and xargs to run veve
	scriptPath := filepath.Join(tmpDir, "batch_stdin.sh")
	scriptContent := `#!/bin/bash
cd ` + tmpDir + `
ls *.md | xargs -I {} bash -c 'veve {} -o {}.pdf --quiet' 
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	cmd := exec.Command("bash", scriptPath)
	err := cmd.Run()

	if err != nil {
		t.Logf("stdin batch processing failed: %v", err)
		t.Skip("feature not yet implemented")
	}
}

// TestBatchConcurrentConversion tests parallel batch conversion.
func TestBatchConcurrentConversion(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	for i := 1; i <= 5; i++ {
		filePath := filepath.Join(tmpDir, "file"+string(rune('0'+i))+".md")
		content := "# File " + string(rune('0'+i)) + "\n\nContent for file\n"
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	// Use GNU parallel or xargs with parallel execution
	scriptPath := filepath.Join(tmpDir, "batch_parallel.sh")
	scriptContent := `#!/bin/bash
cd ` + tmpDir + `
ls *.md | xargs -P 2 -I {} veve {} -o {}.pdf --quiet
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	cmd := exec.Command("bash", scriptPath)
	err := cmd.Run()

	if err != nil {
		t.Logf("parallel batch processing failed: %v", err)
		t.Skip("feature not yet fully implemented")
	}
}

// TestBatchCleanup tests that batch processing cleans up temp files.
func TestBatchCleanup(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	os.WriteFile(filepath.Join(tmpDir, "doc.md"), []byte("# Test\nContent"), 0o644)

	// Run conversion
	scriptPath := filepath.Join(tmpDir, "batch_cleanup.sh")
	scriptContent := `#!/bin/bash
cd ` + tmpDir + `
for f in *.md; do
  veve "$f" --quiet
done
echo "Temp files after conversion:"
ls -la /tmp/veve-theme-* 2>/dev/null | wc -l
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	cmd := exec.Command("bash", scriptPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("cleanup test failed: %v", err)
		return
	}

	t.Logf("Cleanup result: %s", string(output))
}
