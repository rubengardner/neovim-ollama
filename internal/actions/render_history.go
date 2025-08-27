package actions

import (
	"strings"

	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
	"github.com/rubengardner/neovim-ollama/internal/files"
)

func renderHistory(history []files.ChatMessage) string {
	var rendered []string
	styles := ui.NewStyles()
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
