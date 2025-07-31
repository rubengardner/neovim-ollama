package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	input     textinput.Model
	viewport  viewport.Model
	width     int
	height    int
	isWaiting bool
	err       error
	spinner   spinner.Model
	history   []ChatMessage
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
	sp := spinner.New()
	sp.Spinner = spinner.Line
	sp.Style = lipgloss.NewStyle()
	return Model{
		input:    ti,
		viewport: vp,
		spinner:  sp,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
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
		m.viewport.SetContent(renderHistory(m.history))
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.isWaiting {
				return m, nil
			}
			m.isWaiting = true
			m.input.SetValue("")
			m.history = append(m.history, ChatMessage{Role: "user", Content: text}, ChatMessage{Role: "assistant", Content: ""})
			m.viewport.SetContent(renderHistory(m.history))
			cmds = append(cmds, fetchResponse(text, m.history), m.spinner.Tick)
			return m, tea.Batch(cmds...)
		case "up":
			m.viewport.ScrollUp(1)
		case "down":
			m.viewport.ScrollDown(1)
		}

	case responseMsg:
		m.isWaiting = false
		m.input.SetValue("")
		m.input.CursorEnd()
		m.input.Focus()
		m.history[len(m.history)-1].Content = string(msg)
		m.viewport.SetContent(renderHistory(m.history))
		return m, nil

	case errorMsg:
		m.isWaiting = false
		m.input.SetValue("")
		m.input.Focus()
		m.history[len(m.history)-1].Content = "Error: " + msg.Error()
		m.viewport.SetContent(renderHistory(m.history))
		return m, nil
	}

	if m.isWaiting {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)
		m.input.SetValue(m.spinner.View() + " Thinking...")
		m.input.SetCursor(0)
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	m.viewport, _ = m.viewport.Update(msg)
	cmds = append(cmds, inputCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	outputBox := outputStyle.Render(m.viewport.View())
	var spinnerLine string
	inputBox := inputStyle.Render(m.input.View())
	return fmt.Sprintf("%s\n%s%s", outputBox, spinnerLine, inputBox)
}
