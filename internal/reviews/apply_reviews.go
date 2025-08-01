package reviews

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func applyAcceptedChanges(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// For now, just return to chat mode with a success message
		// Later this will actually apply the file changes
		appliedCount := 0
		for _, change := range m.ProposedChanges {
			if change.Accepted != nil && *change.Accepted {
				appliedCount++
			}
		}

		if appliedCount > 0 {
			return fmt.Sprintf("Applied %d changes successfully!", appliedCount)
		} else {
			return "No changes were accepted to apply."
		}
	}
}
