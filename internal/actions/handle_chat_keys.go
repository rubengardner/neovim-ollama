package actions

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/files"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func HandleChatKeys(m model.Model, msg tea.KeyMsg, cmds *[]tea.Cmd) (m.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit

	case "ctrl+f":
		// Switch to file explorer mode
		m.Mode = model.FileExplorerMode
		m.Chat.Input.Placeholder = "Space: select, Enter: open folder, Esc: back to chat"
		*cmds = append(*cmds, loadFiles(m.FileExplorer.CurrentDir))
		return m, tea.Batch(*cmds...)

	case "enter":
		text := strings.TrimSpace(m.Chat.Input.Value())
		if text == "" || m.UI.IsWaiting {
			return m, nil
		}

		// Set waiting state
		m.UI.IsWaiting = true
		m.Chat.Input.SetValue("")

		// Build context from selected files
		contextSummary := buildFileContextSummary(m.FileExplorer.SelectedFiles)
		fullContext := buildFileContext(m.FileExplorer.SelectedFiles)

		// Add user message with summary for display
		userDisplayContent := text
		if contextSummary != "" {
			userDisplayContent = fmt.Sprintf("%s\n\n%s", text, contextSummary)
		}

		// Create full content for AI with file contents
		userAIContent := text
		if fullContext != "" {
			userAIContent = fmt.Sprintf("%s\n\n--- Context from selected files ---\n%s", text, fullContext)
		}

		// Add to history with display version
		m.Chat.History = append(m.Chat.History,
			files.ChatMessage{Role: "user", Content: userDisplayContent},
			files.ChatMessage{Role: "assistant", Content: ""})
		m.Chat.Viewport.SetContent(renderChatHistory(m.Chat.History))

		// Send to AI
		*cmds = append(*cmds, fetchResponseWithContext(userAIContent, m.Chat.History), m.UI.Spinner.Tick)
		return m, tea.Batch(*cmds...)

	case "up":
		m.Chat.Viewport.ScrollUp(1)

	case "down":
		m.Chat.Viewport.ScrollDown(1)
	}

	return m, nil
}
