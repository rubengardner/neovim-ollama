package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/actions"
	"github.com/rubengardner/neovim-ollama/internal/chat"
	"github.com/rubengardner/neovim-ollama/internal/files"
	"github.com/rubengardner/neovim-ollama/internal/reviews"
	"github.com/rubengardner/neovim-ollama/internal/ui"
)

type Mode int

const (
	ChatMode Mode = iota
	FileExplorerMode
	ReviewMode
	FileSelectMode
)

type Model struct {
	Chat         chat.Model
	FileExplorer files.Model
	Review       reviews.Model
	UI           ui.Model
	Mode         Mode
}

func New() Model {
	return Model{
		Chat:         chat.New(),
		FileExplorer: files.New(),
		Review:       reviews.New(),
		UI:           ui.New(),
		Mode:         ChatMode,
	}
}

// Init initializes the model and returns an initial command
func (m Model) Init() tea.Cmd {
	// Return a command batch that initializes all components
	return tea.Batch(
		m.Chat.Init(),
		m.FileExplorer.Init(),
		m.Review.Init(),
		m.UI.Spinner.Tick,
	)
}

// Update handles updates for the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate to the actions.Update function
	return actions.Update(m, msg)
}

// View renders the current state of the model
func (m Model) View() string {
	// Return different views based on the current mode
	switch m.Mode {
	case ChatMode:
		return RenderChatView(m)
	case FileExplorerMode:
		return RenderFileExplorerView(m)
	case ReviewMode:
		return renderReviewView(m)
	default:
		return "Unknown mode"
	}
}
