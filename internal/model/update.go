package model

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		if m.Mode == ChatMode {
			m.Viewport.SetContent(renderHistory(m.History))
		} else if m.Mode == FileSelectMode {
			m.FilesViewport.SetContent(m.renderFiles())
		} else if m.Mode == ReviewMode {
			m.ReviewViewport.SetContent(m.renderReviewChanges())
		}
		return m, nil

	case tea.KeyMsg:
		if m.Mode == ChatMode {
			return m.handleChatKeys(msg, &cmds)
		} else if m.Mode == FileSelectMode {
			return m.handleFileKeys(msg, &cmds)
		} else if m.Mode == ReviewMode {
			return m.handleReviewKeys(msg, &cmds)
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
