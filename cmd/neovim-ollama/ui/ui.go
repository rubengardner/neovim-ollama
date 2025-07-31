package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ollama"
)

type Model struct {
	input    textinput.Model
	viewport viewport.Model
	width    int
	height   int
	history  []string
	err      error
}

type (
	responseMsg string
	errorMsg    error
)

var (
	borderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63"))
	inputStyle    = borderStyle.Padding(0, 1)
	outputStyle   = borderStyle.Padding(0, 1).MarginBottom(1)
	promptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Bold(true)
	responseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F"))
	divider       = lipgloss.NewStyle().Foreground(lipgloss.Color("#444")).Render(strings.Repeat("─", 40))
)

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter prompt"
	ti.Focus()
	ti.CharLimit = 500

	vp := viewport.New(0, 0)

	return Model{
		input:    ti,
		viewport: vp,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = m.width - 4
		m.viewport.Width = m.width - 4
		m.viewport.Height = m.height - 5
		m.viewport.YPosition = 1
		m.viewport.SetContent("")
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				return m, nil
			}
			m.input.SetValue("")
			m.history = append(m.history, promptStyle.Render("You: "+text))
			m.viewport.SetContent(strings.Join(m.history, "\n\n"+divider+"\n\n"))
			return m, fetchResponse(text)
		case "up":
			m.viewport.ScrollUp(1)
		case "down":
			m.viewport.ScrollDown(1)
		}

	case responseMsg:
		rendered := renderMarkdown(string(msg))
		m.history = append(m.history, responseStyle.Render(rendered))
		m.viewport.SetContent(strings.Join(m.history, "\n\n"+divider+"\n\n"))
		return m, nil

	case errorMsg:
		m.err = msg
		return m, nil
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	m.viewport, _ = m.viewport.Update(msg)

	cmds = append(cmds, inputCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	outputBox := outputStyle.Render(m.viewport.View())
	inputBox := inputStyle.Render(m.input.View())
	return fmt.Sprintf("%s\n%s", outputBox, inputBox)
}

func fetchResponse(prompt string) tea.Cmd {
	return func() tea.Msg {
		resp, err := ollama.Generate(prompt)
		if err != nil {
			return errorMsg(err)
		}
		return responseMsg(resp)
	}
}
