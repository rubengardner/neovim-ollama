package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ollama/ollama/api"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ollama"
	"github.com/rubengardner/neovim-ollama/internal/files"
)

func FetchResponseWithContext(fullPrompt string, displayHistory []files.ChatMessage) tea.Cmd {
	return func() tea.Msg {
		messages := []api.Message{
			{Role: "system", Content: "You are a helpful assistant."},
		}

		for i, msg := range displayHistory {
			if (msg.Role == "assistant" && msg.Content == "") ||
				(msg.Role == "user" && i == len(displayHistory)-2) {
				continue
			}

			if msg.Role == "user" && i < len(displayHistory)-2 {
				messages = append(messages, api.Message{
					Role:    msg.Role,
					Content: msg.Content,
				})
			} else if msg.Role == "assistant" {
				messages = append(messages, api.Message{
					Role:    msg.Role,
					Content: msg.Content,
				})
			}
		}

		messages = append(messages, api.Message{
			Role:    "user",
			Content: fullPrompt,
		})

		resp, err := ollama.Generate(messages)
		if err != nil {
			return errorMsg(err)
		}
		return responseMsg(resp)
	}
}
