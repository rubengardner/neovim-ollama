package reviews

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/files"
)

// Model represents the review component for proposed changes
type Model struct {
	ProposedChanges []files.FileChange
	Cursor          int
	Viewport        viewport.Model
}

// New initializes and returns a new review model
func New() Model {
	return Model{
		ProposedChanges: []files.FileChange{},
		Cursor:          0,
		Viewport:        viewport.New(80, 20),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
