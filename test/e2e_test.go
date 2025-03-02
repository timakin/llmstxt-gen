package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2E is an end-to-end test that runs the llmstxt-gen command on test data
// and verifies the output matches the expected result.
func TestE2E(t *testing.T) {
	// Get the project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	rootDir = filepath.Dir(rootDir) // Assuming test is run from the test directory

	var binaryPath string
	if os.Getenv("CI") == "true" {
		// In CI, assume the binary is already built and available
		binaryPath = filepath.Join(rootDir, "llmstxt-gen")
	} else {
		// Build the tool locally
		buildCmd := exec.Command("go", "build", "-o", "llmstxt-gen", ".")
		buildCmd.Dir = rootDir
		buildOutput, err := buildCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to build tool: %v\nOutput: %s", err, buildOutput)
		}
		binaryPath = filepath.Join(rootDir, "llmstxt-gen")
	}

	// Create a temporary output file
	outputFile := filepath.Join(rootDir, "test_output.txt")
	defer os.Remove(outputFile) // Clean up after test

	// Run the tool
	testdataDir := filepath.Join(rootDir, "testdata", "pages")
	cmd := exec.Command(
		binaryPath,
		"--input-dir", testdataDir,
		"--output-file", outputFile,
		"--root-dir", "pages",
		"--project-name", "Test Documentation",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run tool: %v\nOutput: %s", err, output)
	}

	// Read the generated output
	generatedContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Read the expected output
	expectedOutputPath := filepath.Join(rootDir, "testdata", "expected-output.txt")
	expectedContent, err := os.ReadFile(expectedOutputPath)
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	// Compare the output (ignoring whitespace differences)
	generatedStr := normalizeWhitespace(string(generatedContent))
	expectedStr := normalizeWhitespace(string(expectedContent))

	if generatedStr != expectedStr {
		// Find the first difference
		for i := 0; i < len(generatedStr) && i < len(expectedStr); i++ {
			if generatedStr[i] != expectedStr[i] {
				start := max(0, i-20)
				endGen := min(len(generatedStr), i+20)
				endExp := min(len(expectedStr), i+20)

				t.Errorf("Output differs at position %d:\nGenerated: ...%s...\nExpected: ...%s...",
					i, generatedStr[start:endGen], expectedStr[start:endExp])
				break
			}
		}

		if len(generatedStr) != len(expectedStr) {
			t.Errorf("Output length differs: generated=%d, expected=%d",
				len(generatedStr), len(expectedStr))
		}

		t.Fatalf("Generated output does not match expected output")
	}
}

// TestE2EWithOptions tests the tool with different command-line options
func TestE2EWithOptions(t *testing.T) {
	// Get the project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	rootDir = filepath.Dir(rootDir) // Assuming test is run from the test directory

	var binaryPath string
	if os.Getenv("CI") == "true" {
		// In CI, assume the binary is already built and available
		binaryPath = filepath.Join(rootDir, "llmstxt-gen")
	} else {
		// Build the tool locally
		buildCmd := exec.Command("go", "build", "-o", "llmstxt-gen", ".")
		buildCmd.Dir = rootDir
		buildOutput, err := buildCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to build tool: %v\nOutput: %s", err, buildOutput)
		}
		binaryPath = filepath.Join(rootDir, "llmstxt-gen")
	}

	// Create a temporary output file
	outputFile := filepath.Join(rootDir, "test_output_options.txt")
	defer os.Remove(outputFile) // Clean up after test

	// Run the tool with verbose option
	testdataDir := filepath.Join(rootDir, "testdata", "pages")
	cmd := exec.Command(
		binaryPath,
		"--input-dir", testdataDir,
		"--output-file", outputFile,
		"--root-dir", "pages",
		"--project-name", "Custom Project Name",
		"--verbose",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run tool with options: %v\nOutput: %s", err, output)
	}

	// Verify the output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Read the generated output
	generatedContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the custom project name is in the output
	if !strings.Contains(string(generatedContent), "# Custom Project Name") {
		t.Errorf("Custom project name not found in output")
	}

	// Check that verbose output was generated (should be in the command output)
	if !strings.Contains(string(output), "Starting conversion") {
		t.Errorf("Verbose output not found in command output")
	}
}

// normalizeWhitespace removes extra whitespace and normalizes line endings
func normalizeWhitespace(s string) string {
	// Replace all whitespace sequences with a single space
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
