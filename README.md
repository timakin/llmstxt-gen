# llmstxt-gen: MDX to LLMsTXT Converter

This tool converts documentation from MDX format to the LLMsTXT format, which is designed to make documentation more accessible to Large Language Models (LLMs).

## Overview

The LLMsTXT format is a standardized way to provide information to help LLMs use a website at inference time. This tool scans documentation, extracts the content, transforms it into plain Markdown, and formats it according to the LLMsTXT specification.

## Installation

### Option 1: Download pre-built binaries (Recommended)

Download the pre-built binaries from the [Releases](https://github.com/timakin/llmstxt-gen/releases) page.

### Option 2: Install with Homebrew (macOS and Linux)

```bash
# Add the tap
brew tap timakin/tap

# Install the tool
brew install llmstxt-gen
```

### Option 3: Install with go install

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

### Option 4: Clone and Build

```bash
# Clone the repository
git clone https://github.com/timakin/llmstxt-gen.git
cd llmstxt-gen

# Build the tool
go build -o llmstxt-gen
```

## Usage

### If installed with go install

```bash
# Generate llms.txt from documentation
llmstxt-gen --input-dir ./docs --output-file ./llms.txt

# Specify a different root directory for relative path calculation
llmstxt-gen --input-dir ./docs --output-file ./llms.txt --root-dir docs

# Customize the project name
llmstxt-gen --input-dir ./docs --output-file ./llms.txt --project-name "My Project"

# Enable verbose logging
llmstxt-gen --input-dir ./docs --output-file ./llms.txt --verbose
```

### If built from source

```bash
# Generate llms.txt from documentation
./llmstxt-gen --input-dir ./docs --output-file ./llms.txt

# Specify a different root directory for relative path calculation
./llmstxt-gen --input-dir ./docs --output-file ./llms.txt --root-dir docs

# Customize the project name
./llmstxt-gen --input-dir ./docs --output-file ./llms.txt --project-name "My Project"

# Enable verbose logging
./llmstxt-gen --input-dir ./docs --output-file ./llms.txt --verbose
```

Or use the provided shell scripts:

```bash
# Generate llms.txt with default options
./generate-llms.sh

# Generate with custom options
./generate-llms.sh --input-dir ./docs --output-file ./llms.txt --root-dir docs --project-name "My Project" --verbose
```

The shell script will automatically build the tool if it doesn't exist.

### Command-line Options

- `--input-dir`: Input directory containing MDX files (default: "./pages")
- `--output-file`: Output file path (default: "./llms.txt")
- `--root-dir`: Root directory name for relative path calculation (default: "pages")
- `--project-name`: Project name for the LLMsTXT output (default: "Documentation")
- `--verbose`: Enable verbose logging

## How It Works

1. **Scanning**: The tool recursively scans the input directory for MDX files.
2. **Parsing**: Each MDX file is parsed to extract its content, title, and summary.
3. **Transformation**: The MDX content is transformed into plain Markdown, removing React components and other MDX-specific syntax.
4. **Formatting**: The transformed content is formatted according to the LLMsTXT specification.
5. **Output**: The formatted content is written to the output file.

## LLMsTXT Format

The LLMsTXT format includes:

1. An H1 with the name of the project or site (required)
2. A blockquote with a short summary of the project
3. Markdown sections with detailed information
4. Sections delimited by H2 headers containing "file lists" of URLs where further detail is available

For more information about the LLMsTXT format, see [llmstxt.org](https://llmstxt.org/).

## Release Process

This project uses [GoReleaser](https://goreleaser.com/) to automate the release process. Here's how to create a new release:

1. Make sure all your changes are committed and pushed to the repository.
2. Create and push a new tag with the version number:
   ```bash
   git tag -a v0.1.0 -m "First release"
   git push origin v0.1.0
   ```
3. The GitHub Actions workflow will automatically build and publish the release.

You can also test the release process locally without publishing:

```bash
# Install GoReleaser if you haven't already
go install github.com/goreleaser/goreleaser@latest

# Test the release process (dry run)
goreleaser release --snapshot --clean --skip-publish
```

This will create a release in the `dist/` directory without publishing it to GitHub.

## License

This project is open source. Please add an appropriate license file to the repository.
