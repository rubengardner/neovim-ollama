package files

import (
	"fmt"
	"strings"

	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func renderFiles(m *model.Model, styles *ui.Styles) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("Directory: %s\n", m.CurrentDir))
	content.WriteString(fmt.Sprintf("Selected files: %d\n\n", len(m.SelectedFiles)))

	for i, file := range m.Files {
		prefix := "  "
		if i == m.FilesCursor {
			prefix = "> "
		}

		var line string
		if file.IsDir {
			line = styles.Folder.Render(fmt.Sprintf("%s📁 %s/", prefix, file.Name))
		} else {
			fileIcon := "📄"
			if file.Selected {
				fileIcon = "✅"
				line = styles.Checked.Render(fmt.Sprintf("%s%s %s", prefix, fileIcon, file.Name))
			} else {
				line = styles.File.Render(fmt.Sprintf("%s%s %s", prefix, fileIcon, file.Name))
			}
		}

		if i == m.FilesCursor {
			line = styles.Selected.Render(line)
		}

		content.WriteString(line + "\n")
	}

	return content.String()
}
