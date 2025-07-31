package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ollama/ollama/api"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ollama"
)

func fetchResponse(prompt string, history []ChatMessage) tea.Cmd {
	return func() tea.Msg {
		messages := []api.Message{
			{Role: "system", Content: "You are a helpful assistant."},
		}

		for _, msg := range history {
			messages = append(messages, api.Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		messages = append(messages, api.Message{
			Role:    "user",
			Content: prompt,
		})

		resp, err := ollama.Generate(messages)
		if err != nil {
			return errorMsg(err)
		}
		return responseMsg(resp)
	}
}
