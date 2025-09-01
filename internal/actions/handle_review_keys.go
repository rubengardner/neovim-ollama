package actions

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/files"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func HandleReviewKeys(m model.Model, msg tea.KeyMsg, cmds *[]tea.Cmd) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.Mode = model.ChatMode
		return m, nil

	case "up", "k":
		if m.Review.Cursor > 0 {
			m.Review.Cursor--
			m.Review.Viewport.SetContent(renderReviewChanges(m.Review.ProposedChanges))
		}

	case "down", "j":
		if m.Review.Cursor < len(m.Review.ProposedChanges)-1 {
			m.Review.Cursor++
			m.Review.Viewport.SetContent(renderReviewChanges(m.Review.ProposedChanges))
		}

	case "enter":
		if len(m.Review.ProposedChanges) > 0 {
			// Logic to apply the change would go here
			// ...
		}
	}

	return m, nil
}

func renderReviewChanges(changes []files.FileChange) string {
	var content strings.Builder

	content.WriteString("Proposed changes:\n\n")

	for i, change := range changes {
		marker := "  "
		if i == 0 {
			marker = "> "
		}

		content.WriteString(fmt.Sprintf("%s%s (lines %d-%d)\n",
			marker,
			change.FilePath,
			change.LineStart,
			change.LineEnd))
	}

	return content.String()
}
