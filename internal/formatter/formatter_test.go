package formatter

import (
	"strings"
	"testing"

	"github.com/basemachina/llmstxt-gen/internal/parser"
)

func TestFormatLLMsTXT(t *testing.T) {
	// Create test data
	contents := []parser.ParsedContent{
		{
			FilePath:     "section1/test.mdx",
			Title:        "Test Document",
			Summary:      "This is a test document",
			Content:      "# Test Document\n\nThis is the content of the test document.",
			RelativePath: "section1/test.mdx",
			Section:      "section1",
		},
		{
			FilePath:     "section2/another-test.mdx",
			Title:        "Another Test Document",
			Summary:      "This is another test document",
			Content:      "# Another Test Document\n\nThis is the content of another test document.",
			RelativePath: "section2/another-test.mdx",
			Section:      "section2",
		},
	}

	// Format the content
	projectName := "Test Project"
	result := FormatLLMsTXT(contents, projectName)

	// Verify the result
	if !strings.Contains(result, "# Test Project") {
		t.Errorf("Project name not included in the output: %s", result)
	}

	if !strings.Contains(result, "> Test Project is a documentation site.") {
		t.Errorf("Summary not included in the output: %s", result)
	}

	if !strings.Contains(result, "## Section1") {
		t.Errorf("Section1 header not included in the output: %s", result)
	}

	if !strings.Contains(result, "## Section2") {
		t.Errorf("Section2 header not included in the output: %s", result)
	}

	if !strings.Contains(result, "- [Test Document](/section1/test)") {
		t.Errorf("File list for section1 not included in the output: %s", result)
	}

	if !strings.Contains(result, "- [Another Test Document](/section2/another-test)") {
		t.Errorf("File list for section2 not included in the output: %s", result)
	}

	if !strings.Contains(result, "### Test Document") {
		t.Errorf("Content for Test Document not included in the output: %s", result)
	}

	if !strings.Contains(result, "### Another Test Document") {
		t.Errorf("Content for Another Test Document not included in the output: %s", result)
	}
}

func TestFormatLLMsTXTWithOptions(t *testing.T) {
	// Create test data
	contents := []parser.ParsedContent{
		{
			FilePath:     "section1/test.mdx",
			Title:        "Test Document",
			Summary:      "This is a test document",
			Content:      "# Test Document\n\nThis is the content of the test document.",
			RelativePath: "section1/test.mdx",
			Section:      "section1",
		},
	}

	// Create custom options
	options := FormatOptions{
		ProjectName:      "Custom Project",
		Summary:          "Custom summary for the project",
		GeneralInfo:      "Custom general information",
		OrganizationInfo: "Custom organization information",
	}

	// Format the content with custom options
	result := FormatLLMsTXTWithOptions(contents, options)

	// Verify the result
	if !strings.Contains(result, "# Custom Project") {
		t.Errorf("Custom project name not included in the output: %s", result)
	}

	if !strings.Contains(result, "> Custom summary for the project") {
		t.Errorf("Custom summary not included in the output: %s", result)
	}

	if !strings.Contains(result, "Custom general information") {
		t.Errorf("Custom general information not included in the output: %s", result)
	}

	if !strings.Contains(result, "Custom organization information") {
		t.Errorf("Custom organization information not included in the output: %s", result)
	}
}

func TestDefaultFormatOptions(t *testing.T) {
	projectName := "Test Project"
	options := DefaultFormatOptions(projectName)

	if options.ProjectName != projectName {
		t.Errorf("Expected project name %s, got %s", projectName, options.ProjectName)
	}

	if !strings.Contains(options.Summary, projectName) {
		t.Errorf("Project name not included in the summary: %s", options.Summary)
	}

	if !strings.Contains(options.GeneralInfo, projectName) {
		t.Errorf("Project name not included in the general info: %s", options.GeneralInfo)
	}

	if options.OrganizationInfo != "The documentation is organized by topic." {
		t.Errorf("Unexpected organization info: %s", options.OrganizationInfo)
	}
}

func TestGroupBySection(t *testing.T) {
	// Create test data
	contents := []parser.ParsedContent{
		{
			FilePath:     "section1/test1.mdx",
			Title:        "Test Document 1",
			Summary:      "This is test document 1",
			Content:      "Content 1",
			RelativePath: "section1/test1.mdx",
			Section:      "section1",
		},
		{
			FilePath:     "section1/test2.mdx",
			Title:        "Test Document 2",
			Summary:      "This is test document 2",
			Content:      "Content 2",
			RelativePath: "section1/test2.mdx",
			Section:      "section1",
		},
		{
			FilePath:     "section2/test3.mdx",
			Title:        "Test Document 3",
			Summary:      "This is test document 3",
			Content:      "Content 3",
			RelativePath: "section2/test3.mdx",
			Section:      "section2",
		},
	}

	// Group by section
	sectionMap := groupBySection(contents)

	// Verify the result
	if len(sectionMap) != 2 {
		t.Errorf("Expected 2 sections, got %d", len(sectionMap))
	}

	if len(sectionMap["section1"]) != 2 {
		t.Errorf("Expected 2 documents in section1, got %d", len(sectionMap["section1"]))
	}

	if len(sectionMap["section2"]) != 1 {
		t.Errorf("Expected 1 document in section2, got %d", len(sectionMap["section2"]))
	}

	if sectionMap["section1"][0].Title != "Test Document 1" {
		t.Errorf("Expected first document in section1 to be 'Test Document 1', got '%s'", sectionMap["section1"][0].Title)
	}

	if sectionMap["section1"][1].Title != "Test Document 2" {
		t.Errorf("Expected second document in section1 to be 'Test Document 2', got '%s'", sectionMap["section1"][1].Title)
	}

	if sectionMap["section2"][0].Title != "Test Document 3" {
		t.Errorf("Expected first document in section2 to be 'Test Document 3', got '%s'", sectionMap["section2"][0].Title)
	}
}

func TestFormatSectionTitle(t *testing.T) {
	tests := []struct {
		name    string
		section string
		want    string
	}{
		{
			name:    "Special case: action",
			section: "action",
			want:    "Actions",
		},
		{
			name:    "Special case: view",
			section: "view",
			want:    "Views",
		},
		{
			name:    "Special case: admin",
			section: "admin",
			want:    "Administration",
		},
		{
			name:    "Special case: faq",
			section: "faq",
			want:    "FAQ",
		},
		{
			name:    "Special case: tips",
			section: "tips",
			want:    "Tips and Tricks",
		},
		{
			name:    "Regular case: single word",
			section: "section",
			want:    "Section",
		},
		{
			name:    "Regular case: multiple words with underscore",
			section: "getting_started",
			want:    "Getting Started",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSectionTitle(tt.section)
			if got != tt.want {
				t.Errorf("formatSectionTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
