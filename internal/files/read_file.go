package files

import (
	"fmt"
	"path/filepath"
	"strings"
)

func BuildFileContext(SelectedFiles []string) string {
	if len(SelectedFiles) == 0 {
		return ""
	}

	var context strings.Builder
	for _, filePath := range SelectedFiles {
		content, err := readFileContent(filePath)
		if err != nil {
			context.WriteString(fmt.Sprintf("Error reading %s: %s\n\n", filePath, err.Error()))
			continue
		}

		optimizedContent := optimizeFileContent(content, filePath)

		context.WriteString(fmt.Sprintf("File: %s\n", filePath))
		context.WriteString("```\n")
		context.WriteString(optimizedContent)
		context.WriteString("\n```\n\n")
	}

	return context.String()
}

func BuildFileContextSummary(SelectedFiles []string) string {
	if len(SelectedFiles) == 0 {
		return ""
	}

	var summary strings.Builder
	summary.WriteString("📁 Attached files:\n")
	for _, filePath := range SelectedFiles {
		fileName := filepath.Base(filePath)
		relPath, err := filepath.Rel(".", filePath)
		if err != nil {
			relPath = filePath
		}
		summary.WriteString(fmt.Sprintf("  • %s (%s)\n", fileName, relPath))
	}

	return summary.String()
}

func optimizeFileContent(content, filePath string) string {
	lines := strings.Split(content, "\n")
	var optimized []string

	emptyLineCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			emptyLineCount++
			if emptyLineCount <= 1 {
				optimized = append(optimized, "")
			}
			continue
		}
		emptyLineCount = 0

		if shouldStripComments(filePath) {
			if strings.HasPrefix(trimmed, "//") ||
				strings.HasPrefix(trimmed, "#") ||
				strings.HasPrefix(trimmed, "/*") ||
				strings.HasPrefix(trimmed, "*") {
				continue
			}
		}

		optimized = append(optimized, line)
	}

	result := strings.Join(optimized, "\n")

	const maxSize = 8000
	if len(result) > maxSize {
		truncated := result[:maxSize-200]
		result = truncated + "\n\n... (file truncated for brevity) ..."
	}

	return result
}

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
