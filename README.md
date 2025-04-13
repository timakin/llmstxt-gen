# llmstxt-gen: HTML to LLMsTXT Converter

[![CI](https://github.com/timakin/llmstxt-gen/actions/workflows/ci.yml/badge.svg)](https://github.com/timakin/llmstxt-gen/actions/workflows/ci.yml)
[![Release](https://github.com/timakin/llmstxt-gen/actions/workflows/release.yml/badge.svg)](https://github.com/timakin/llmstxt-gen/actions/workflows/release.yml)

This tool converts HTML files into the LLMsTXT format, designed to make website content more accessible to Large Language Models (LLMs). It extracts readable content from HTML using the `go-readability` library.

## Overview

The LLMsTXT format is a standardized way to provide information to help LLMs understand and utilize website content effectively. This tool processes HTML files, either by scanning a directory or using a sitemap, extracts the main textual content, and formats it according to the LLMsTXT specification.

## Installation

### Option 1: Download pre-built binaries (Recommended)

Download the pre-built binaries from the [Releases](https://github.com/timakin/llmstxt-gen/releases) page.

### Option 2: Install with go install

#### From GitHub

```bash
# Install the latest version
go install github.com/timakin/llmstxt-gen@latest

# The binary will be installed to your $GOPATH/bin directory
# Make sure $GOPATH/bin is in your PATH
```

#### From local repository

```bash
# Clone the repository
git clone git@github.com:timakin/llmstxt-gen.git
cd llmstxt-gen

# Install the tool
go install .

# The binary will be installed to your $GOPATH/bin directory
# Make sure $GOPATH/bin is in your PATH
```

### Option 3: Clone and Build

```bash
# Clone the repository
git clone https://github.com/timakin/llmstxt-gen.git
cd llmstxt-gen

# Build the tool
go build -o llmstxt-gen
```

## Usage

### Basic Usage (Scanning a Directory)

```bash
# Generate llms.txt from HTML files in ./public directory
llmstxt-gen --html-dir ./public --output-file ./llms.txt

# Customize the project name
llmstxt-gen --html-dir ./public --output-file ./llms.txt --project-name "My Website"

# Enable verbose logging
llmstxt-gen --html-dir ./public --output-file ./llms.txt --verbose
```

### Using a Sitemap

```bash
# Generate llms.txt using a sitemap.xml located in ./public
# Assumes HTML files are also in ./public and sitemap URLs map accordingly
llmstxt-gen --html-dir ./public --sitemap ./public/sitemap.xml --output-file ./llms.txt

# Specify a different project name with sitemap
llmstxt-gen --html-dir ./public --sitemap ./public/sitemap.xml --output-file ./llms.txt --project-name "My Blog"
```

### Command-line Options

- `--html-dir`: Input directory containing HTML files (default: "./html"). This directory is scanned if `--sitemap` is not provided. It's also used to find local files corresponding to sitemap URLs.
- `--sitemap`: Path to the sitemap XML file (optional). If provided, only URLs listed in the sitemap will be processed.
- `--output-file`: Output file path (default: "./llms.txt").
- `--project-name`: Project name for the LLMsTXT output (default: "Documentation").
- `--verbose`: Enable verbose logging.
- `--version`, `-v`: Display version information.

## How It Works

1.  **Input Source Determination**: Checks if a `--sitemap` path is provided.
2.  **File List Generation**:
    *   **Sitemap Mode**: Parses the sitemap, extracts URLs, and attempts to map each URL to a corresponding local HTML file within the `--html-dir`.
    *   **Directory Scan Mode**: Recursively scans the `--html-dir` for `.html` and `.htm` files.
3.  **Content Extraction**: For each identified HTML file:
    *   Opens the file.
    *   Uses `go-readability` (`github.com/mackee/go-readability`) to extract the main readable content (title, plain text content, excerpt).
4.  **Formatting**: Organizes the extracted content (title, URL, excerpt, full text) into sections based on the directory structure relative to `--html-dir`. Formats the collected information according to the LLMsTXT specification.
5.  **Output**: Writes the formatted content to the specified `--output-file`.

## LLMsTXT Format

The LLMsTXT format includes:

1.  An H1 with the name of the project or site (required).
2.  A blockquote with a short summary of the project.
3.  Markdown sections with detailed information, often grouped by topic or directory structure.
4.  Sections delimited by H2 headers containing file lists (links to specific pages/documents) with summaries.
5.  Detailed content for each page under H3 headers.

For more information about the LLMsTXT format, see [llmstxt.org](https://llmstxt.org/).

## Release Process

This project uses [GoReleaser](https://goreleaser.com/) to automate the release process. Here's how to create a new release:

1.  Make sure all your changes are committed and pushed to the repository.
2.  Create and push a new tag with the version number:
    ```bash
    git tag -a vX.Y.Z -m "Release vX.Y.Z"
    git push origin vX.Y.Z
    ```
3.  The GitHub Actions workflow will automatically build and publish the release.

You can also test the release process locally without publishing:

```bash
# Install GoReleaser if you haven't already
# go install github.com/goreleaser/goreleaser@latest

# Test the release process (dry run)
goreleaser release --snapshot --clean --skip=publish
```

This will create a release in the `dist/` directory without publishing it to GitHub.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
