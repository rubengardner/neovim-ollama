package actions

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
	"github.com/rubengardner/neovim-ollama/internal/files"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func HandleChatKeys(m model.Model, msg tea.KeyMsg, cmds *[]tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "ctrl+f":
		m.Mode = model.FileSelectMode
		m.Input.Placeholder = "Space: select, Enter: open folder, Esc: back to chat"
		*cmds = append(*cmds, files.LoadFiles(m.CurrentDir))
		return m, tea.Batch(*cmds...)
	case "enter":
		text := strings.TrimSpace(m.Input.Value())
		if text == "" || m.IsWaiting {
			return m, nil
		}

		m.IsWaiting = true
		m.Input.SetValue("")

		contextSummary := files.BuildFileContextSummary(&m)
		fullContext := files.BuildFileContext(&m)

		userDisplayContent := text
		if contextSummary != "" {
			userDisplayContent = fmt.Sprintf("%s\n\n%s", text, contextSummary)
		}

		userAIContent := text
		if fullContext != "" {
			userAIContent = fmt.Sprintf("%s\n\n--- Context from selected files ---\n%s", text, fullContext)
		}

		m.History = append(m.History,
			files.ChatMessage{Role: "user", Content: userDisplayContent},
			files.ChatMessage{Role: "assistant", Content: ""})
		m.Viewport.SetContent(ui.RenderHistory(m))

		*cmds = append(*cmds, ui.FetchResponseWithContext(userAIContent, m.History), m.Spinner.Tick)
		return m, tea.Batch(*cmds...)
	case "up":
		m.Viewport.ScrollUp(1)
	case "down":
		m.Viewport.ScrollDown(1)
	default:
		if !m.IsWaiting {
			var inputCmd tea.Cmd
			m.Input, inputCmd = m.Input.Update(msg)
			if inputCmd != nil {
				*cmds = append(*cmds, inputCmd)
			}
		}
	}
	return m, nil
}
