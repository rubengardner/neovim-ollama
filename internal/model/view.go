package model

import (
	"fmt"
)

func (m Model) View() string {
	if m.Mode == FileSelectMode {
		outputBox := m.Styles.Output.Render(m.FilesViewport.View())
		inputBox := m.Styles.Input.Render(m.Input.View())
		helpText := m.Styles.Help.Render("↑/↓: navigate, Space: select file, Enter: open folder, Esc: back to chat")
		return fmt.Sprintf("%s\n%s\n%s", outputBox, inputBox, helpText)
	}

	if m.Mode == ReviewMode {
		outputBox := m.Styles.Output.Render(m.ReviewViewport.View())
		inputBox := m.Styles.Input.Render(m.Input.View())
		helpText := m.Styles.Help.Render("↑/↓: navigate, Space: accept, r: reject, u: undo, Enter: apply all, Esc: back to chat")
		return fmt.Sprintf("%s\n%s\n%s", outputBox, inputBox, helpText)
	}

	// Chat mode
	outputBox := m.Styles.Output.Render(m.Viewport.View())
	var spinnerLine string
	inputBox := m.Styles.Input.Render(m.Input.View())

	var helpText string
	if len(m.SelectedFiles) > 0 {
		helpText = m.Styles.Help.Render(fmt.Sprintf("Context: %d files selected | Ctrl+F: file browser, Ctrl+R: review mode", len(m.SelectedFiles)))
	} else {
		helpText = m.Styles.Help.Render("Ctrl+F: file browser, Ctrl+R: review mode, ↑/↓: scroll")
	}

	return fmt.Sprintf("%s\n%s%s\n%s", outputBox, spinnerLine, inputBox, helpText)
}
