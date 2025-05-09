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

	// Always build the tool to ensure it's available
	binaryPath := filepath.Join(rootDir, "llmstxt-gen")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = rootDir
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build tool: %v\nOutput: %s", err, buildOutput)
	}

	// Verify the binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("Binary was not created at %s", binaryPath)
	}

	// Create a temporary output file
	outputFile := filepath.Join(rootDir, "test_output.txt")
	defer os.Remove(outputFile) // Clean up after test

	// Run the tool
	testdataDir := filepath.Join(rootDir, "testdata", "html") // Use new html testdata dir
	cmd := exec.Command(
		binaryPath,
		"--html-dir", testdataDir, // Use --html-dir
		"--output-file", outputFile,
		// Remove --root-dir
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

	// Always build the tool to ensure it's available
	binaryPath := filepath.Join(rootDir, "llmstxt-gen")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = rootDir
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build tool: %v\nOutput: %s", err, buildOutput)
	}

	// Verify the binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("Binary was not created at %s", binaryPath)
	}

	// Create a temporary output file
	outputFile := filepath.Join(rootDir, "test_output_options.txt")
	defer os.Remove(outputFile) // Clean up after test

	// Run the tool with verbose option
	testdataDir := filepath.Join(rootDir, "testdata", "html") // Use new html testdata dir
	cmd := exec.Command(
		binaryPath,
		"--html-dir", testdataDir, // Use --html-dir
		"--output-file", outputFile,
		// Remove --root-dir
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

// TestE2EWithSitemap tests the tool with sitemap option
func TestE2EWithSitemap(t *testing.T) {
	// Get the project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	rootDir = filepath.Dir(rootDir) // Assuming test is run from the test directory

	// Always build the tool to ensure it's available
	binaryPath := filepath.Join(rootDir, "llmstxt-gen")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = rootDir
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build tool: %v\nOutput: %s", err, buildOutput)
	}

	// Verify the binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("Binary was not created at %s", binaryPath)
	}

	// Create a temporary output file
	outputFile := filepath.Join(rootDir, "test_output_sitemap.txt")
	defer os.Remove(outputFile) // Clean up after test

	// Run the tool with sitemap option
	testdataDir := filepath.Join(rootDir, "testdata", "html")
	sitemapPath := filepath.Join(rootDir, "testdata", "sitemap.xml")
	cmd := exec.Command(
		binaryPath,
		"--html-dir", testdataDir,
		"--sitemap", sitemapPath,
		"--output-file", outputFile,
		"--project-name", "Test Documentation",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run tool with sitemap: %v\nOutput: %s", err, output)
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
