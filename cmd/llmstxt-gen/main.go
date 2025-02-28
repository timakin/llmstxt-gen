package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/basemachina/llmstxt-gen/internal/formatter"
	"github.com/basemachina/llmstxt-gen/internal/parser"
	"github.com/basemachina/llmstxt-gen/internal/transformer"
)

var (
	inputDir    = flag.String("input-dir", "./pages", "Input directory containing MDX files")
	outputFile  = flag.String("output-file", "./llms.txt", "Output file path")
	rootDir     = flag.String("root-dir", "pages", "Root directory name for relative path calculation")
	projectName = flag.String("project-name", "Documentation", "Project name for the LLMsTXT output")
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
	flag.Parse()

	// Validate input directory
	info, err := os.Stat(*inputDir)
	if err != nil {
		log.Fatalf("Error accessing input directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("Input path is not a directory: %s", *inputDir)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(*outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	if *verbose {
		log.Printf("Starting conversion from %s to %s", *inputDir, *outputFile)
	}

	// Scan for MDX files
	files, err := scanMDXFiles(*inputDir)
	if err != nil {
		log.Fatalf("Error scanning MDX files: %v", err)
	}

	if *verbose {
		log.Printf("Found %d MDX files", len(files))
	}

	// Parse and transform content
	var contents []parser.ParsedContent
	for _, file := range files {
		if *verbose {
			log.Printf("Processing file: %s", file)
		}

		// Parse MDX file
		content, err := parser.ParseMDXFile(file, *rootDir)
		if err != nil {
			log.Printf("Error parsing file %s: %v", file, err)
			continue
		}

		// Transform content
		transformed := transformer.TransformContent(content)
		contents = append(contents, transformed)
	}

	// Format content according to LLMsTXT specification
	llmsTxtContent := formatter.FormatLLMsTXT(contents, *projectName)

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

// scanMDXFiles recursively scans the input directory for MDX files
func scanMDXFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file is an MDX file
		if filepath.Ext(path) == ".mdx" || filepath.Ext(path) == ".md" {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
