package chat

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderChatHistory(history []ChatMessage) string {
	var content strings.Builder

	for _, msg := range history {
		switch msg.Role {
		case "user":
			promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Bold(true)
			content.WriteString(promptStyle.Render("You: ") + "\n")
			content.WriteString(msg.Content + "\n\n")
		case "assistant":
			responseStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F"))
			content.WriteString(responseStyle.Render("Assistant: ") + "\n")
			content.WriteString(msg.Content + "\n\n")
		}
	}

	return content.String()
}
