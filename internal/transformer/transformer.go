package transformer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/timakin/llmstxt-gen/internal/parser"
)

// TransformContent transforms the parsed MDX content into plain Markdown
func TransformContent(content parser.ParsedContent) parser.ParsedContent {
	// Create a copy of the content
	transformed := content

	// Transform the content
	transformedContent := content.Content

	// First, protect code blocks by replacing them with placeholders
	codeBlocks := make(map[string]string)
	transformedContent = protectCodeBlocks(transformedContent, codeBlocks)

	// Remove import statements
	transformedContent = removeImports(transformedContent)

	// Replace MDX components with plain Markdown
	transformedContent = replaceComponents(transformedContent)

	// Replace JSX expressions
	transformedContent = replaceJSXExpressions(transformedContent)

	// Clean up any remaining MDX-specific syntax
	transformedContent = cleanupMDX(transformedContent)

	// Restore code blocks
	transformedContent = restoreCodeBlocks(transformedContent, codeBlocks)

	// Update the transformed content
	transformed.Content = transformedContent

	return transformed
}

// protectCodeBlocks replaces code blocks with placeholders to protect them from other transformations
func protectCodeBlocks(content string, codeBlocks map[string]string) string {
	// Regular expression to match code blocks (```language\n...```)
	codeBlockRegex := regexp.MustCompile("(?ms)```([a-zA-Z0-9]*)\n(.*?)```")

	// Replace each code block with a placeholder
	result := codeBlockRegex.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := fmt.Sprintf("CODE_BLOCK_%d", len(codeBlocks))
		codeBlocks[placeholder] = match
		return placeholder
	})

	return result
}

// restoreCodeBlocks replaces placeholders with the original code blocks
func restoreCodeBlocks(content string, codeBlocks map[string]string) string {
	result := content

	// Replace each placeholder with the original code block
	for placeholder, codeBlock := range codeBlocks {
		result = strings.Replace(result, placeholder, codeBlock, 1)
	}

	return result
}

// removeImports removes import statements from the content
func removeImports(content string) string {
	importRegex := regexp.MustCompile(`(?m)^import.*from.*$`)
	result := importRegex.ReplaceAllString(content, "")

	// Clean up multiple consecutive newlines that might be left after removing imports
	multipleNewlinesRegex := regexp.MustCompile(`\n{3,}`)
	result = multipleNewlinesRegex.ReplaceAllString(result, "\n\n")

	return result
}

// replaceComponents replaces MDX components with plain Markdown
func replaceComponents(content string) string {
	// Replace common components
	replacements := map[string]string{
		// Information component
		`<Information>([^<]*)</Information>`: "> **Information:** $1",

		// DocImage component
		`<DocImage\s+alt="([^"]*)"\s+src="([^"]*)"\s*/>`: "![Image: $1]($2)",

		// iframe (e.g., YouTube videos)
		`<iframe[^>]*src="([^"]*)"[^>]*></iframe>`: "> **Video:** [$1]($1)",

		// Other common components
		`<Flexbox[^>]*>`: "",
		`</Flexbox>`:     "",
		`<Box[^>]*>`:     "",
		`</Box>`:         "",
		`<Card[^>]*>`:    "",
		`</Card>`:        "",
		`<Grid[^>]*>`:    "",
		`</Grid>`:        "",
		`<Heading[^>]*text=\{?["']([^"']*)["']\}?[^>]*>`: "### $1\n",
		`</Heading>`: "",
		`<Button[^>]*title=\{?["']([^"']*)["']\}?[^>]*>`: "[$1]",
		`</Button>`: "",
	}

	result := content
	for pattern, replacement := range replacements {
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, replacement)
	}

	// Handle custom MDX components that span multiple lines
	// This is a simplified approach; a proper parser would be more robust
	// Note: Go's regexp doesn't support backreferences, so we'll use a simpler approach
	// Look for common closing tags
	for _, tag := range []string{"Information", "DocImage", "Card", "Grid", "Box", "Flexbox", "Button", "Heading"} {
		openTag := "<" + tag + "[^>]*>"
		closeTag := "</" + tag + ">"
		re := regexp.MustCompile(openTag + `(?s)(.*?)` + closeTag)
		result = re.ReplaceAllString(result, "$1")
	}

	// Handle self-closing custom components
	selfClosingComponentRegex := regexp.MustCompile(`<([A-Z][a-zA-Z]*)[^>/]*/>`)
	result = selfClosingComponentRegex.ReplaceAllString(result, "")

	return result
}

// replaceJSXExpressions replaces JSX expressions with plain text
func replaceJSXExpressions(content string) string {
	// Replace {variable} expressions
	jsxExprRegex := regexp.MustCompile(`\{([^{}]*)\}`)
	return jsxExprRegex.ReplaceAllString(content, "[Expression: $1]")
}

// cleanupMDX cleans up any remaining MDX-specific syntax
func cleanupMDX(content string) string {
	// Remove export statements
	exportRegex := regexp.MustCompile(`(?m)^export default.*$`)
	content = exportRegex.ReplaceAllString(content, "")

	// Remove empty lines at the beginning and end
	content = strings.TrimSpace(content)

	// Replace multiple consecutive empty lines with a single empty line
	multipleEmptyLinesRegex := regexp.MustCompile(`\n{3,}`)
	content = multipleEmptyLinesRegex.ReplaceAllString(content, "\n\n")

	return content
}
