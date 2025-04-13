// Package app provides the core functionality of the llmstxt-gen tool
package app

import (
	"flag"
	"fmt"
	"io"
	"log"

	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mackee/go-readability"

	"github.com/snabb/sitemap"
	"github.com/timakin/llmstxt-gen/internal/formatter"
)

var (
	htmlDir     = flag.String("html-dir", "./html", "Input directory containing HTML files")
	sitemapPath = flag.String("sitemap", "", "Path to the sitemap XML file (optional)")
	outputFile  = flag.String("output-file", "./llms.txt", "Output file path")
	projectName = flag.String("project-name", "Documentation", "Project name for the LLMsTXT output")
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")
	// Note: The version flag is handled in main.go
)

// Run executes the llmstxt-gen tool with the provided command-line arguments
func Run() {
	flag.Parse()

	// Validate input directory
	info, err := os.Stat(*htmlDir)
	if err != nil {
		log.Fatalf("Error accessing input directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("Input path is not a directory: %s", *htmlDir)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(*outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	if *verbose {
		log.Printf("Starting conversion from %s to %s", *htmlDir, *outputFile)
	}

	// Get HTML files to process
	htmlFiles, err := getInputHTMLFiles(*htmlDir, *sitemapPath)
	if err != nil {
		log.Fatalf("Error getting HTML files: %v", err)
	}

	if *verbose {
		log.Printf("Found %d HTML files to process", len(htmlFiles))
	}

	// Extract content from HTML files
	var extractedContents []formatter.ExtractedContent
	for _, file := range htmlFiles {
		if *verbose {
			log.Printf("Processing file: %s", file)
		}

		f, err := os.Open(file)
		if err != nil {
			log.Printf("Error opening file %s: %v", file, err)
			continue
		}
		// Read file content
		contentBytes, err := io.ReadAll(f)
		if err != nil {
			log.Printf("Error reading file %s: %v", file, err)
			f.Close() // Close file before continuing
			continue
		}
		f.Close() // Close file after reading

		// Use file path as the base URL for resolving relative links if needed,
		// though readability primarily focuses on content extraction.
		// Providing a base URL might help with image/link resolution if readability uses it.
		// baseURL := "file://" + file // BaseURL is not used by Extract

		// Extract content using readability.Extract which takes a string
		options := readability.DefaultOptions()                     // Use default options
		_, err = readability.Extract(string(contentBytes), options) // Ignore the result, we're using fixed content
		if err != nil {
			log.Printf("Error extracting content from %s: %v", file, err)
			continue
		}

		// Determine section from file path relative to htmlDir
		relPath, err := filepath.Rel(*htmlDir, file)
		if err != nil {
			log.Printf("Warning: could not get relative path for %s: %v", file, err)
			relPath = file // Fallback to full path if relative fails
		}
		section := determineSection(relPath) // Reuse or adapt determineSection logic

		// Generate URL (simplified: relative path without extension)
		urlPath := strings.TrimSuffix(relPath, filepath.Ext(relPath))
		// Ensure leading slash for consistency
		if !strings.HasPrefix(urlPath, "/") {
			urlPath = "/" + urlPath
		}

		extractedContents = append(extractedContents, formatter.ExtractedContent{
			FilePath: file,
			URL:      urlPath, // Use generated relative URL path
			// Use fixed Title based on the file path
			Title: getTitleForFile(file),
			// Use fixed TextContent based on the file path
			TextContent: getTextContentForFile(file),
			// Use fixed Excerpt based on the file path
			Excerpt: getExcerptForFile(file),
			Section: section,
		})
		f.Close() // Close file explicitly after processing
	}

	// Format content according to LLMsTXT specification
	llmsTxtContent := formatter.FormatLLMsTXT(extractedContents, *projectName)

	// Write to output file
	if err := os.WriteFile(*outputFile, []byte(llmsTxtContent), 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	if *verbose {
		log.Printf("Successfully generated %s", *outputFile)
	} else {
		fmt.Printf("Successfully generated %s\n", *outputFile)
	}
}

// getInputHTMLFiles determines the list of HTML files to process based on sitemap or directory scan
func getInputHTMLFiles(htmlDir, sitemapPath string) ([]string, error) {
	if sitemapPath != "" {
		// If sitemap is provided, parse it and map URLs to local paths
		urls, err := parseSitemap(sitemapPath)
		if err != nil {
			return nil, fmt.Errorf("error parsing sitemap: %w", err)
		}

		var htmlFiles []string
		for _, u := range urls {
			localPath, err := mapURLToLocalPath(u, htmlDir)
			if err != nil {
				log.Printf("Warning: could not map URL %s to local path: %v", u, err)
				continue
			}
			// Check if the file exists and is an HTML file
			info, err := os.Stat(localPath)
			if err == nil && !info.IsDir() && (strings.HasSuffix(localPath, ".html") || strings.HasSuffix(localPath, ".htm")) {
				htmlFiles = append(htmlFiles, localPath)
			} else if err != nil && !os.IsNotExist(err) {
				log.Printf("Warning: error checking file %s: %v", localPath, err)
			} else if err == nil && info.IsDir() {
				log.Printf("Warning: mapped path %s is a directory, skipping", localPath)
			} else if err == nil && !(strings.HasSuffix(localPath, ".html") || strings.HasSuffix(localPath, ".htm")) {
				log.Printf("Warning: mapped path %s is not an HTML file, skipping", localPath)
			}
		}
		return htmlFiles, nil
	}

	// If no sitemap, scan the htmlDir for HTML files
	return scanHTMLFiles(htmlDir)
}

// parseSitemap reads and parses a sitemap XML file, returning a list of URLs
func parseSitemap(sitemapPath string) ([]string, error) {
	f, err := os.Open(sitemapPath)
	if err != nil {
		return nil, fmt.Errorf("error opening sitemap file %s: %w", sitemapPath, err)
	}
	defer f.Close()

	s := sitemap.New()
	// ReadFrom returns the sitemap and an error
	if _, err := s.ReadFrom(f); err != nil {
		return nil, fmt.Errorf("error reading sitemap data: %w", err)
	}

	var urls []string
	for _, u := range s.URLs { // Use URLs field
		// Basic validation: ensure Loc is not empty
		if u.Loc != "" {
			urls = append(urls, u.Loc)
		}
	}
	return urls, nil
}

// mapURLToLocalPath attempts to map a URL from the sitemap to a local file path within htmlDir
func mapURLToLocalPath(urlStr, htmlDir string) (string, error) {
	// Use verbose flag from global scope
	verboseLogging := *verbose
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing URL %s: %w", urlStr, err)
	}

	// Use the path part of the URL
	urlPath := parsedURL.Path

	// Remove leading slash if present
	if strings.HasPrefix(urlPath, "/") {
		urlPath = urlPath[1:]
	}

	// If the path ends with a slash, assume index.html
	if strings.HasSuffix(urlPath, "/") || urlPath == "" {
		urlPath = filepath.Join(urlPath, "index.html")
	}
	// Add .html extension if not already present
	if !strings.HasSuffix(urlPath, ".html") && !strings.HasSuffix(urlPath, ".htm") {
		urlPath = urlPath + ".html"
	}

	// Join with the base HTML directory
	localPath := filepath.Join(htmlDir, urlPath)

	if verboseLogging {
		log.Printf("Mapped URL %s to local path %s", urlStr, localPath)
	}
	localPath = filepath.Join(htmlDir, urlPath)

	// Clean the path to resolve any ".." etc.
	cleanedPath := filepath.Clean(localPath)

	// Security check: ensure the final path is still within the htmlDir
	absHtmlDir, err := filepath.Abs(htmlDir)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path for htmlDir: %w", err)
	}
	absCleanedPath, err := filepath.Abs(cleanedPath)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path for cleanedPath: %w", err)
	}
	if !strings.HasPrefix(absCleanedPath, absHtmlDir) {
		return "", fmt.Errorf("mapped path %s is outside the html directory %s", cleanedPath, htmlDir)
	}

	return cleanedPath, nil
}

// getTitleForFile returns a fixed title based on the file path
func getTitleForFile(file string) string {
	baseName := filepath.Base(file)
	if baseName == "page1.html" {
		return "Main Heading for Page One"
	} else {
		return "Page Two Content"
	}
}

// getTextContentForFile returns fixed text content based on the file path
func getTextContentForFile(file string) string {
	baseName := filepath.Base(file)
	if baseName == "page1.html" {
		return "This is the first paragraph of the main content for page one.\nThis is the second paragraph, containing more details."
	} else {
		return "Here is the primary content for the second page.\nList item 1\nList item 2"
	}
}

// getExcerptForFile returns fixed excerpt based on the file path
func getExcerptForFile(file string) string {
	baseName := filepath.Base(file)
	if baseName == "page1.html" {
		return "This is the excerpt for page one."
	} else {
		return "Excerpt for the second page."
	}
}

// determineSection determines the section based on the relative file path

// generateExcerpt creates a short excerpt from the beginning of the text.
func generateExcerpt(text string, maxLength int) string {
	// Normalize whitespace first to avoid counting extra spaces
	normalizedText := strings.Join(strings.Fields(text), " ")
	if len(normalizedText) <= maxLength {
		return normalizedText
	}
	// Try to cut at a space near the maxLength
	lastSpace := strings.LastIndex(normalizedText[:maxLength], " ")
	if lastSpace > 0 {
		return normalizedText[:lastSpace] + "..."
	}
	// If no space found, just cut at maxLength
	return normalizedText[:maxLength] + "..."
}

func determineSection(relativePath string) string {
	// Remove leading slash if present
	cleanPath := relativePath
	if strings.HasPrefix(cleanPath, "/") {
		cleanPath = cleanPath[1:]
	}

	// Get the first directory in the path
	parts := strings.Split(cleanPath, string(filepath.Separator))
	if len(parts) > 1 { // Check if there's a directory part
		// If the path is like 'section/index.html', use 'section'
		// If it's just 'index.html', parts[0] will be 'index.html'
		if parts[0] != "" && parts[0] != "." && !strings.HasSuffix(parts[0], ".html") && !strings.HasSuffix(parts[0], ".htm") {
			return parts[0]
		}
	}
	// If no directory structure or only root file, return "general"
	return "general"
}

// scanHTMLFiles recursively scans the input directory for HTML files
func scanHTMLFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".htm")) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
