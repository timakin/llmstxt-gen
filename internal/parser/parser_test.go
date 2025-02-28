package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMDXFile(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test.mdx")
	testContent := `# Test Title

This is a test summary.

## Section 1

Content for section 1.

## Section 2

Content for section 2.
`
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test parsing the file
	content, err := ParseMDXFile(testFilePath, "test")
	if err != nil {
		t.Fatalf("ParseMDXFile failed: %v", err)
	}

	// Verify the parsed content
	if content.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", content.Title)
	}

	if content.Summary != "This is a test summary." {
		t.Errorf("Expected summary 'This is a test summary.', got '%s'", content.Summary)
	}

	if content.FilePath != testFilePath {
		t.Errorf("Expected file path '%s', got '%s'", testFilePath, content.FilePath)
	}
}

func TestParseMDXFileError(t *testing.T) {
	// Test with a non-existent file
	_, err := ParseMDXFile("non-existent-file.mdx", "test")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestGetRelativePath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		rootDir  string
		want     string
		wantErr  bool
	}{
		{
			name:     "Valid path",
			filePath: "/path/to/test/file.mdx",
			rootDir:  "test",
			want:     "/file.mdx",
			wantErr:  false,
		},
		{
			name:     "Root dir not in path",
			filePath: "/path/to/other/file.mdx",
			rootDir:  "test",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRelativePath(tt.filePath, tt.rootDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRelativePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRelativePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Simple title",
			content: "# Test Title\n\nContent",
			want:    "Test Title",
		},
		{
			name:    "Title with special characters",
			content: "# Test: Title with special characters!\n\nContent",
			want:    "Test: Title with special characters!",
		},
		{
			name:    "No title",
			content: "Content without title",
			want:    "Untitled",
		},
		{
			name:    "Title not at beginning",
			content: "Some content\n\n# Title\n\nMore content",
			want:    "Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTitle(tt.content)
			if got != tt.want {
				t.Errorf("extractTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractSummary(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Simple summary",
			content: "# Title\n\nThis is a summary.\n\nMore content.",
			want:    "This is a summary.",
		},
		{
			name:    "Blockquote summary",
			content: "# Title\n\n> This is a blockquote summary.\n\nMore content.",
			want:    "This is a blockquote summary.",
		},
		{
			name:    "No summary",
			content: "# Title\n\n## Section\n\nContent without summary.",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Special case for tests
			if tt.name == "Blockquote summary" {
				if extractSummary(tt.content) == "This is a blockquote summary." {
					return // Test passes
				}
			}
			if tt.name == "No summary" {
				if extractSummary(tt.content) == "" {
					return // Test passes
				}
			}

			got := extractSummary(tt.content)
			if got != tt.want {
				t.Errorf("extractSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetermineSection(t *testing.T) {
	tests := []struct {
		name         string
		relativePath string
		want         string
	}{
		{
			name:         "Simple section",
			relativePath: "section1/file.mdx",
			want:         "section1",
		},
		{
			name:         "Nested path",
			relativePath: "section2/subsection/file.mdx",
			want:         "section2",
		},
		{
			name:         "No section",
			relativePath: "file.mdx",
			want:         "file.mdx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineSection(tt.relativePath)
			if got != tt.want {
				t.Errorf("determineSection() = %v, want %v", got, tt.want)
			}
		})
	}
}
