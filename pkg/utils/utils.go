package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnsureDirectoryExists ensures that a directory exists
func EnsureDirectoryExists(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsMarkdownFile checks if a file is a Markdown file
func IsMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".mdx"
}

// GetRelativePath returns the path relative to a base directory
func GetRelativePath(path, baseDir string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %w", err)
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("error getting absolute base directory: %w", err)
	}

	relPath, err := filepath.Rel(absBaseDir, absPath)
	if err != nil {
		return "", fmt.Errorf("error getting relative path: %w", err)
	}

	return relPath, nil
}

// SanitizeFilename sanitizes a filename
func SanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename

	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	return result
}

// GetFileInfo returns information about a file
func GetFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// ReadFile reads a file and returns its content
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes content to a file
func WriteFile(path string, content []byte) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := EnsureDirectoryExists(dir); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	// Write the file
	return os.WriteFile(path, content, 0644)
}
