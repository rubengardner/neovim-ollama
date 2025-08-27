package actions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
	"github.com/rubengardner/neovim-ollama/internal/model"
	"github.com/rubengardner/neovim-ollama/internal/reviews"
)

func Update(m model.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Input.Width = m.Width - 4
		m.Viewport.Width = m.Width - 4
		m.Viewport.Height = m.Height - 7
		m.Viewport.YPosition = 1
		m.FilesViewport.Width = m.Width - 4
		m.FilesViewport.Height = m.Height - 7
		m.FilesViewport.YPosition = 1
		m.ReviewViewport.Width = m.Width - 4
		m.ReviewViewport.Height = m.Height - 7
		m.ReviewViewport.YPosition = 1

		if m.Mode == model.ChatMode {
			m.Viewport.SetContent(ui.RenderHistory(m.History))
		} else if m.Mode == model.FileSelectMode {
			m.FilesViewport.SetContent(ui.RenderFiles(m))
		} else if m.Mode == model.ReviewMode {
			m.ReviewViewport.SetContent(reviews.RenderReviewChanges())
		}
		return m, nil

	case tea.KeyMsg:
		if m.Mode == model.ChatMode {
			return HandleChatKeys(m, msg, &cmds)
		} else if m.Mode == models.FileSelectMode {
			return HandleFileKeys(msg, &cmds)
		} else if m.Mode == ReviewMode {
		}

	case responseMsg:
		m.IsWaiting = false
		m.Input.SetValue("")
		m.Input.CursorEnd()
		m.Input.Focus()
		if len(m.History) > 0 {
			m.History[len(m.History)-1].Content = string(msg)
			m.Viewport.SetContent(renderHistory(m.History))
		}
		return m, nil

	case errorMsg:
		m.IsWaiting = false
		m.Input.SetValue("")
		m.Input.Focus()
		if len(m.History) > 0 {
			m.History[len(m.History)-1].Content = "Error: " + msg.Error()
			m.Viewport.SetContent(renderHistory(m.History))
		}
		return m, nil

	case filesLoadedMsg:
		m.Files = []FileItem(msg)
		m.FilesCursor = 0
		m.FilesViewport.SetContent(m.renderFiles())
		return m, nil

	case reviewChangesMsg:
		m.ProposedChanges = []FileChange(msg)
		m.ReviewCursor = 0
		m.Mode = ReviewMode
		m.Input.Placeholder = "Space: accept, r: reject, Enter: apply all, Esc: back to chat"
		m.ReviewViewport.SetContent(m.renderReviewChanges())
		return m, nil
	}

	if m.IsWaiting {
		var spinCmd tea.Cmd
		m.Spinner, spinCmd = m.Spinner.Update(msg)
		cmds = append(cmds, spinCmd)
		m.Input.SetValue(m.Spinner.View() + " Thinking...")
		m.Input.SetCursor(0)
	} else {
		if m.Mode != ChatMode {
			var inputCmd tea.Cmd
			m.Input, inputCmd = m.Input.Update(msg)
			cmds = append(cmds, inputCmd)
		}
	}

	if m.Mode == ChatMode {
		m.Viewport, _ = m.Viewport.Update(msg)
	} else if m.Mode == FileSelectMode {
		m.FilesViewport, _ = m.FilesViewport.Update(msg)
	} else if m.Mode == ReviewMode {
		m.ReviewViewport, _ = m.ReviewViewport.Update(msg)
	}

	return m, tea.Batch(cmds...)
}
