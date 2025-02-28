package transformer

import (
	"strings"
	"testing"

	"github.com/basemachina/llmstxt-gen/internal/parser"
)

func TestTransformContent(t *testing.T) {
	// Create a test ParsedContent
	input := parser.ParsedContent{
		FilePath:     "test.mdx",
		Title:        "Test Title",
		Summary:      "Test Summary",
		Content:      "# Test Title\n\n<Information>This is an information component.</Information>\n\n```go\nfunc main() {}\n```\n\nexport default () => <div>Test</div>",
		RelativePath: "test.mdx",
		Section:      "test",
	}

	// Transform the content
	result := TransformContent(input)

	// Verify the transformation
	if result.Title != input.Title {
		t.Errorf("Expected title '%s', got '%s'", input.Title, result.Title)
	}

	if result.Summary != input.Summary {
		t.Errorf("Expected summary '%s', got '%s'", input.Summary, result.Summary)
	}

	// Check that the Information component was transformed
	if !strings.Contains(result.Content, "> **Information:** This is an information component.") {
		t.Errorf("Information component not transformed correctly: %s", result.Content)
	}

	// Check that the code block was preserved
	if !strings.Contains(result.Content, "```go\nfunc main() {}\n```") {
		t.Errorf("Code block not preserved: %s", result.Content)
	}

	// Check that the export statement was removed
	if strings.Contains(result.Content, "export default") {
		t.Errorf("Export statement not removed: %s", result.Content)
	}
}

func TestProtectCodeBlocks(t *testing.T) {
	content := "Before code\n\n```go\nfunc main() {}\n```\n\nAfter code\n\n```js\nconsole.log('test');\n```"
	codeBlocks := make(map[string]string)

	result := protectCodeBlocks(content, codeBlocks)

	// Check that code blocks were replaced with placeholders
	if strings.Contains(result, "```go") || strings.Contains(result, "```js") {
		t.Errorf("Code blocks not replaced with placeholders: %s", result)
	}

	// Check that we have 2 code blocks in the map
	if len(codeBlocks) != 2 {
		t.Errorf("Expected 2 code blocks, got %d", len(codeBlocks))
	}

	// Restore code blocks
	restored := restoreCodeBlocks(result, codeBlocks)

	// Check that the content was restored correctly
	if restored != content {
		t.Errorf("Content not restored correctly.\nExpected: %s\nGot: %s", content, restored)
	}
}

func TestRemoveImports(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Simple import",
			content: "import React from 'react';\n\nContent",
			want:    "\n\nContent",
		},
		{
			name:    "Multiple imports",
			content: "import React from 'react';\nimport { Component } from 'react';\n\nContent",
			want:    "\n\nContent",
		},
		{
			name:    "No imports",
			content: "Content without imports",
			want:    "Content without imports",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeImports(tt.content)
			if got != tt.want {
				t.Errorf("removeImports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceComponents(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		check   func(string, string) bool
	}{
		{
			name:    "Information component",
			content: "<Information>This is info</Information>",
			want:    "> **Information:** This is info",
			check:   strings.Contains,
		},
		{
			name:    "DocImage component",
			content: "<DocImage alt=\"Test\" src=\"/test.png\" />",
			want:    "![Image: Test](/test.png)",
			check:   strings.Contains,
		},
		{
			name:    "Iframe component",
			content: "<iframe src=\"https://example.com\" width=\"560\" height=\"315\"></iframe>",
			want:    "> **Video:** [https://example.com](https://example.com)",
			check:   strings.Contains,
		},
		{
			name:    "Nested components",
			content: "<Card><Heading text=\"Title\">Content</Heading><Button title=\"Click\"/></Card>",
			want:    "### Title\nContent[Click]",
			check: func(got, want string) bool {
				// Simplified check for nested components
				return !strings.Contains(got, "<Card>") &&
					!strings.Contains(got, "<Heading") &&
					!strings.Contains(got, "<Button") &&
					strings.Contains(got, "Title") &&
					strings.Contains(got, "Content") &&
					strings.Contains(got, "Click")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceComponents(tt.content)
			if !tt.check(got, tt.want) {
				t.Errorf("replaceComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceJSXExpressions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Simple expression",
			content: "The value is {value}.",
			want:    "The value is [Expression: value].",
		},
		{
			name:    "Expression with operation",
			content: "The result is {result + 1}.",
			want:    "The result is [Expression: result + 1].",
		},
		{
			name:    "No expressions",
			content: "Content without expressions",
			want:    "Content without expressions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceJSXExpressions(tt.content)
			if got != tt.want {
				t.Errorf("replaceJSXExpressions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanupMDX(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Export statement",
			content: "Content\n\nexport default () => <div>Test</div>",
			want:    "Content",
		},
		{
			name:    "Multiple empty lines",
			content: "Line 1\n\n\n\nLine 2",
			want:    "Line 1\n\nLine 2",
		},
		{
			name:    "Whitespace at beginning and end",
			content: "\n\nContent\n\n",
			want:    "Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanupMDX(tt.content)
			if got != tt.want {
				t.Errorf("cleanupMDX() = %v, want %v", got, tt.want)
			}
		})
	}
}
