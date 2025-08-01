package files

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (m Model) buildFileContext() string {
	if len(m.selectedFiles) == 0 {
		return ""
	}

	var context strings.Builder
	for _, filePath := range m.selectedFiles {
		content, err := readFileContent(filePath)
		if err != nil {
			context.WriteString(fmt.Sprintf("Error reading %s: %s\n\n", filePath, err.Error()))
			continue
		}

		// Compress/optimize content
		optimizedContent := optimizeFileContent(content, filePath)

		context.WriteString(fmt.Sprintf("File: %s\n", filePath))
		context.WriteString("```\n")
		context.WriteString(optimizedContent)
		context.WriteString("\n```\n\n")
	}

	return context.String()
}

func (m Model) buildFileContextSummary() string {
	if len(m.selectedFiles) == 0 {
		return ""
	}

	var summary strings.Builder
	summary.WriteString("📁 Attached files:\n")
	for _, filePath := range m.selectedFiles {
		fileName := filepath.Base(filePath)
		// Show relative path if it's in current working directory
		relPath, err := filepath.Rel(".", filePath)
		if err != nil {
			relPath = filePath
		}
		summary.WriteString(fmt.Sprintf("  • %s (%s)\n", fileName, relPath))
	}

	return summary.String()
}

// Optimize file content for AI consumption
func optimizeFileContent(content, filePath string) string {
	// Remove excessive whitespace and empty lines
	lines := strings.Split(content, "\n")
	var optimized []string

	emptyLineCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip excessive empty lines (keep max 1 consecutive empty line)
		if trimmed == "" {
			emptyLineCount++
			if emptyLineCount <= 1 {
				optimized = append(optimized, "")
			}
			continue
		}
		emptyLineCount = 0

		// Remove comments for certain file types to reduce size (optional)
		if shouldStripComments(filePath) {
			if strings.HasPrefix(trimmed, "//") ||
				strings.HasPrefix(trimmed, "#") ||
				strings.HasPrefix(trimmed, "/*") ||
				strings.HasPrefix(trimmed, "*") {
				continue
			}
		}

		optimized = append(optimized, line) // Keep original indentation
	}

	result := strings.Join(optimized, "\n")

	// If still too large, truncate with message
	const maxSize = 8000 // Adjust based on your needs
	if len(result) > maxSize {
		truncated := result[:maxSize-200] // Leave room for message
		result = truncated + "\n\n... (file truncated for brevity) ..."
	}

	return result
}

// Decide whether to strip comments based on file type
func shouldStripComments(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	stripCommentExts := []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".c"}

	for _, stripExt := range stripCommentExts {
		if ext == stripExt {
			return true
		}
	}
	return false
}
