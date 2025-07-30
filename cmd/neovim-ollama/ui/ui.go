package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ollama"
)

type Model struct {
	input        textinput.Model
	viewport     viewport.Model
	width        int
	height       int
	isStreaming  bool
	err          error
	streamCancel context.CancelFunc
	history      []string
}

type (
	responseChunkMsg string
	errorMsg         error
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
			if m.isStreaming && m.streamCancel != nil {
				m.streamCancel()
			}
			return m, tea.Quit

		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.isStreaming {
				return m, nil
			}

			ctx, cancel := context.WithCancel(context.Background())
			m.streamCancel = cancel
			m.input.SetValue("")
			m.viewport.SetContent("")
			m.history = append(m.history, promptStyle.Render("You: "+text), "Thinking...")
			m.viewport.SetContent(strings.Join(m.history, "\n\n"+divider+"\n\n"))
			m.isStreaming = true
			return m, streamResponse(ctx, text)
		case "up":
			m.viewport.ScrollUp(1)
		case "down":
			m.viewport.ScrollDown(1)
		}

	case responseChunkMsg:
		rendered, err := glamour.Render(string(msg), "dark")
		if err != nil {
			m.viewport.SetContent(string(msg))
		} else {
			m.history[len(m.history)-1] = responseStyle.Render(rendered)
			m.viewport.SetContent(strings.Join(m.history, "\n\n"+divider+"\n\n"))
		}
		m.isStreaming = false
		return m, nil

	case errorMsg:
		m.err = msg
		m.isStreaming = false
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

func streamResponse(ctx context.Context, prompt string) tea.Cmd {
	return func() tea.Msg {
		var builder strings.Builder

		err := ollama.StreamGenerate(prompt, func(chunk ollama.StreamChunk) {
			select {
			case <-ctx.Done():
				return
			default:
				builder.WriteString(chunk.Response)
			}
		})
		if err != nil {
			return errorMsg(err)
		}
		return responseChunkMsg(builder.String())
	}
}
