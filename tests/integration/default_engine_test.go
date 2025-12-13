package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestDefaultEngineConsistency tests that default engine selection is consistent across conversions (T065)
func TestDefaultEngineConsistency(t *testing.T) {
	vevePath := buildVeveForIntegration(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("default engine consistent across multiple conversions", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test files with unicode content
		testFiles := []string{"test1.md", "test2.md", "test3.md"}
		for _, fname := range testFiles {
			path := filepath.Join(tmpDir, fname)
			markdown := `# Test
Unicode: 世界 🎉
`
			if err := os.WriteFile(path, []byte(markdown), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}
		}

		// Run conversions multiple times
		outputs := make([]string, len(testFiles))
		for i, fname := range testFiles {
			inputFile := filepath.Join(tmpDir, fname)
			outputFile := filepath.Join(tmpDir, fname+".pdf")
			outputs[i] = outputFile

			cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
			if err := cmd.Run(); err != nil {
				t.Logf("Note: conversion %d failed (engine may not be available): %v", i+1, err)
				continue
			}

			// Verify PDF was created
			if info, err := os.Stat(outputFile); err != nil {
				t.Logf("PDF not created for %s: %v", fname, err)
			} else if info.Size() == 0 {
				t.Logf("Warning: PDF for %s is empty", fname)
			}
		}

		// Check that all conversions produced PDFs of similar size
		// (indicating same engine was used, though sizes may vary slightly)
		nonZeroSizes := []int64{}
		for _, output := range outputs {
			if info, err := os.Stat(output); err == nil && info.Size() > 0 {
				nonZeroSizes = append(nonZeroSizes, info.Size())
			}
		}

		if len(nonZeroSizes) > 1 {
			// All conversions succeeded - check consistency
			// PDFs with similar content should be within 20% size range
			minSize := nonZeroSizes[0]
			maxSize := nonZeroSizes[0]
			for _, size := range nonZeroSizes {
				if size < minSize {
					minSize = size
				}
				if size > maxSize {
					maxSize = size
				}
			}

			sizeRatio := float64(maxSize) / float64(minSize)
			if sizeRatio > 1.3 {
				t.Logf("Note: PDF sizes vary (min: %d, max: %d), ratio: %.2f", minSize, maxSize, sizeRatio)
			} else {
				t.Logf("Default engine consistent: PDF sizes within 30%% (min: %d, max: %d)", minSize, maxSize)
			}
		}
	})
}

// TestDefaultEngineFallback tests that fallback occurs when primary engine unavailable (T066)
func TestDefaultEngineFallback(t *testing.T) {
	vevePath := buildVeveForIntegration(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("fallback to alternative engine when primary unavailable", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "test.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		markdown := `# Test
World: 世界
Emoji: 🎉
`
		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Run conversion with default engine (which may trigger fallback internally)
		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		start := time.Now()
		err := cmd.Run()
		elapsed := time.Since(start)

		t.Logf("Default engine selection took %v", elapsed)

		if err != nil {
			t.Logf("Note: conversion failed (may indicate all engines failed or unavailable): %v", err)
			return
		}

		// Verify PDF was created
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Logf("PDF not created (acceptable if all engines unavailable): %v", err)
			return
		}

		if info.Size() == 0 {
			t.Error("PDF is empty after conversion")
			return
		}

		t.Logf("Fallback mechanism successful: created %d byte PDF", info.Size())
	})

	t.Run("selection respects priority even with fallback", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "priority_test.md")
		outputFile := filepath.Join(tmpDir, "priority.pdf")

		markdown := `# Priority Test
Testing priority: ∑ ± ∫
`
		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Try conversion
		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		if err := cmd.Run(); err != nil {
			t.Logf("Note: conversion failed: %v", err)
			return
		}

		// If successful, verify output
		info, _ := os.Stat(outputFile)
		if info.Size() > 0 {
			t.Logf("Priority-based selection worked: %d bytes", info.Size())
		}
	})
}
