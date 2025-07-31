package ui

import "strings"

type ChatMessage struct {
	Role    string
	Content string
}

func renderHistory(history []ChatMessage) string {
	var rendered []string
	for _, msg := range history {
		switch msg.Role {
		case "user":
			rendered = append(rendered, promptStyle.Render("You: "+msg.Content))
		case "assistant":
			rendered = append(rendered, responseStyle.Render(renderMarkdown(msg.Content)))
		default:
			rendered = append(rendered, msg.Content)
		}
	}
	return strings.Join(rendered, "\n\n"+divider+"\n\n")
}
