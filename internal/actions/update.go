package actions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
	"github.com/rubengardner/neovim-ollama/internal/chat"
	"github.com/rubengardner/neovim-ollama/internal/files"
	"github.com/rubengardner/neovim-ollama/internal/model"
	"github.com/rubengardner/neovim-ollama/internal/reviews"
)

func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UI.Width = msg.Width
		m.UI.Height = msg.Height

		cmd := updateViewportDimensions(&m, msg.Width, msg.Height)
		return m, cmd

	case tea.KeyMsg:
		switch m.Mode {
		case model.ChatMode:
			return HandleChatKeys(m, msg, &cmds)
		case model.FileExplorerMode:
			return HandleFileKeys(m, msg, &cmds)
		case model.ReviewMode:
			return HandleReviewKeys(m, msg, &cmds)
		}

	case string:
		return handleResponse(m, string(msg))

	case error:
		return handleError(m, msg)

	case []files.FileItem:
		return handleFilesLoaded(m, msg)

	case []files.FileChange:
		return reviews.HandleReviewChanges(m, msg)
	}

	// Handle spinner updates when waiting
	if m.UI.IsWaiting {
		var spinCmd tea.Cmd
		m.UI.Spinner, spinCmd = m.UI.Spinner.Update(msg)
		cmds = append(cmds, spinCmd)
	}

	// Update
	switch m.Mode {
	case model.ChatMode:
		m.Chat.Viewport, _ = m.Chat.Viewport.Update(msg)
	case model.FileExplorerMode:
		m.FileExplorer.Viewport, _ = m.FileExplorer.Viewport.Update(msg)
	case model.ReviewMode:
		m.Review.Viewport, _ = m.Review.Viewport.Update(msg)
	}

	if !m.UI.IsWaiting {
		var inputCmd tea.Cmd
		m.Chat.Input, inputCmd = m.Chat.Input.Update(msg)
		cmds = append(cmds, inputCmd)
	}

	return m, tea.Batch(cmds...)
}

func updateViewportDimensions(m *model.Model, width, height int) tea.Cmd {
	contentWidth := width - 4
	contentHeight := height - 7

	// Update chat viewport
	m.Chat.Viewport.Width = contentWidth
	m.Chat.Viewport.Height = contentHeight
	m.Chat.Viewport.YPosition = 1
	m.Chat.Input.Width = contentWidth

	// Update file explorer viewport
	m.FileExplorer.Viewport.Width = contentWidth
	m.FileExplorer.Viewport.Height = contentHeight
	m.FileExplorer.Viewport.YPosition = 1

	// Update review viewport
	m.Review.Viewport.Width = contentWidth
	m.Review.Viewport.Height = contentHeight
	m.Review.Viewport.YPosition = 1

	switch m.Mode {
	case model.ChatMode:
		m.Chat.Viewport.SetContent(ui.RenderHistory(m.Chat.History))
	case model.FileExplorerMode:
		m.FileExplorer.Viewport.SetContent(files.RenderFileList(m.FileExplorer))
	case model.ReviewMode:
		m.Review.Viewport.SetContent(renderReviewChanges(m.Review.ProposedChanges))
	}

	return nil
}

func handleFilesLoaded(m model.Model, msg []files.FileItem) (model.Model, tea.Cmd) {
	m.FileExplorer.Files = msg
	m.FileExplorer.Cursor = 0
	m.FileExplorer.Viewport.SetContent(files.RenderFileList(m.FileExplorer))
	return m, nil
}

// Handle response from AI
func handleResponse(m model.Model, response string) (model.Model, tea.Cmd) {
	m.UI.IsWaiting = false
	m.Chat.Input.SetValue("")
	m.Chat.Input.Focus()

	m.Chat.History[len(m.Chat.History)-1].Content = response
	m.Chat.Viewport.SetContent(chat.RenderChatHistory(m.Chat.History))
	return m, nil
}

// Handle error
func handleError(m model.Model, err error) (model.Model, tea.Cmd) {
	m.UI.IsWaiting = false
	m.UI.Err = err
	m.Chat.Input.SetValue("")
	m.Chat.Input.Focus()

	if len(m.Chat.History) > 0 {
		m.Chat.History[len(m.Chat.History)-1].Content = "Error: " + err.Error()
		m.Chat.Viewport.SetContent(chat.RenderChatHistory(m.Chat.History))
	}

	return m, nil
}
