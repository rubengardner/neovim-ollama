package actions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/chat"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func HandleFileKeys(m model.Model, msg tea.KeyMsg, cmds *[]tea.Cmd) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		// Switch back to chat mode
		m.Mode = model.ChatMode
		m.Chat.Input.Placeholder = "Enter prompt (Ctrl+F for files, Ctrl+C to exit)"
		m.Chat.Viewport.SetContent(chat.RenderChatHistory(m.Chat.History))
		return m, nil

	case "up", "k":
		if m.FileExplorer.Cursor > 0 {
			m.FileExplorer.Cursor--
			m.FileExplorer.Viewport.SetContent(renderFileList(m.FileExplorer))
		}

	case "down", "j":
		if m.FileExplorer.Cursor < len(m.FileExplorer.Files)-1 {
			m.FileExplorer.Cursor++
			m.FileExplorer.Viewport.SetContent(renderFileList(m.FileExplorer))
		}

	case " ": // Space to toggle selection
		if len(m.FileExplorer.Files) > 0 && !m.FileExplorer.Files[m.FileExplorer.Cursor].IsDir {
			filePath := m.FileExplorer.Files[m.FileExplorer.Cursor].Path

			// Toggle selection
			found := false
			for i, selected := range m.FileExplorer.SelectedFiles {
				if selected == filePath {
					m.FileExplorer.SelectedFiles = append(m.FileExplorer.SelectedFiles[:i], m.FileExplorer.SelectedFiles[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				m.FileExplorer.SelectedFiles = append(m.FileExplorer.SelectedFiles, filePath)
			}

			// Update selected state in file list
			m.FileExplorer.Files[m.FileExplorer.Cursor].Selected = !found
			m.FileExplorer.Viewport.SetContent(renderFileList(m.FileExplorer))
		}

	case "enter":
		if len(m.FileExplorer.Files) > 0 && m.FileExplorer.Files[m.FileExplorer.Cursor].IsDir {
			m.FileExplorer.CurrentDir = m.FileExplorer.Files[m.FileExplorer.Cursor].Path
			*cmds = append(*cmds, loadFiles(m.FileExplorer.CurrentDir))
			return m, tea.Batch(*cmds...)
		}
	}

	return m, nil
}
