package reviews

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/files"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

type reviewChangesMsg []files.FileChange

func HandleReviewChanges(m model.Model, changes reviewChangesMsg) (model.Model, tea.Cmd) {
	m.Review.ProposedChanges = []files.FileChange(changes)
	m.Review.Cursor = 0
	m.Review.Viewport.SetContent(RenderReviewChanges(m.Review.ProposedChanges))
	return m, nil
}
