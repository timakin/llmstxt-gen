#!/bin/bash

# Script to generate llms.txt from documentation using llmstxt-gen

# Set default values
INPUT_DIR="./pages"
OUTPUT_FILE="./llms.txt"
ROOT_DIR="pages"
PROJECT_NAME="Documentation"
VERBOSE=false

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --input-dir)
      INPUT_DIR="$2"
      shift 2
      ;;
    --output-file)
      OUTPUT_FILE="$2"
      shift 2
      ;;
    --root-dir)
      ROOT_DIR="$2"
      shift 2
      ;;
    --project-name)
      PROJECT_NAME="$2"
      shift 2
      ;;
    --verbose)
      VERBOSE=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Build the tool if it doesn't exist
if [ ! -f "./llmstxt-gen" ]; then
  echo "Building llmstxt-gen tool..."
  go build -o llmstxt-gen ./cmd/llmstxt-gen
fi

# Run the tool
echo "Generating llms.txt..."
if [ "$VERBOSE" = true ]; then
  ./llmstxt-gen --input-dir "$INPUT_DIR" --output-file "$OUTPUT_FILE" --root-dir "$ROOT_DIR" --project-name "$PROJECT_NAME" --verbose
else
  ./llmstxt-gen --input-dir "$INPUT_DIR" --output-file "$OUTPUT_FILE" --root-dir "$ROOT_DIR" --project-name "$PROJECT_NAME"
fi

# Check if the generation was successful
if [ $? -eq 0 ]; then
  echo "Successfully generated $OUTPUT_FILE"
  echo "File size: $(du -h "$OUTPUT_FILE" | cut -f1)"
else
  echo "Error generating $OUTPUT_FILE"
  exit 1
fi
