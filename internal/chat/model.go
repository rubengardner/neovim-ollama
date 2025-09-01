package chat

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Input    textinput.Model
	Viewport viewport.Model
	History  []ChatMessage
	IsTyping bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "Ask a question..."
	input.Focus()

	return Model{
		Input:    input,
		Viewport: viewport.New(20, 80),
		History:  []ChatMessage{},
	}
}

type ChatMessage struct {
	Role    string
	Content string
}
