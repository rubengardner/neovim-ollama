package ui

import (
	"strings"

	"github.com/rubengardner/neovim-ollama/internal/files"
)

func RenderHistory(chatMessage []files.ChatMessage) string {
	var rendered []string
	for _, msg := range chatMessage {
		switch msg.Role {
		case "user":
			rendered = append(rendered, promptStyle.Render("You: "+msg.Content))
		case "assistant":
			rendered = append(rendered, responseStyle.Render(RenderMarkdown(msg.Content)))
		default:
			rendered = append(rendered, msg.Content)
		}
	}
	return strings.Join(rendered, "\n\n"+divider+"\n\n")
}
