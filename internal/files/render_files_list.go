package files

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderFileList(fe Model) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("Directory: %s\n", fe.CurrentDir))
	content.WriteString(fmt.Sprintf("Selected files: %d\n\n", len(fe.SelectedFiles)))

	for i, file := range fe.Files {
		prefix := "  "
		if i == fe.Cursor {
			prefix = "> "
		}

		var line string
		if file.IsDir {
			folderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FABD2F")).Bold(true)
			line = folderStyle.Render(fmt.Sprintf("%s◯ %s/", prefix, file.Name))
		} else {
			fileIcon := "◯"
			if file.Selected {
				fileIcon = "✓"
				checkedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).Bold(true)
				line = checkedStyle.Render(fmt.Sprintf("%s%s %s", prefix, fileIcon, file.Name))
			} else {
				fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#83A598"))
				line = fileStyle.Render(fmt.Sprintf("%s%s %s", prefix, fileIcon, file.Name))
			}
		}

		if i == fe.Cursor {
			selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("#3C3836")).Foreground(lipgloss.Color("#EBDBB2"))
			line = selectedStyle.Render(line)
		}

		content.WriteString(line + "\n")
	}

	return content.String()
}
