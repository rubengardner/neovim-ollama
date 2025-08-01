package model

import (
	"strings"

	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
)

func renderHistory(history []ChatMessage, styles *ui.Styles) string {
	var rendered []string
	for _, msg := range history {
		switch msg.Role {
		case "user":
			rendered = append(rendered, styles.Prompt.Render("You: "+msg.Content))
		case "assistant":
			rendered = append(rendered, styles.Response.Render(ui.RenderMarkdown(msg.Content)))
		default:
			rendered = append(rendered, msg.Content)
		}
	}
	return strings.Join(rendered, "\n\n"+styles.Divider+"\n\n")
}
