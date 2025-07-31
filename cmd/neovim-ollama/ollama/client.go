package ollama

import (
	"context"

	"github.com/ollama/ollama/api"
)

func Generate(messages []api.Message) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return "", err
	}

	stream := false
	req := &api.ChatRequest{
		Model:    "llama3",
		Messages: messages,
		Stream:   &stream,
	}

	var fullResponse string
	ctx := context.Background()
	err = client.Chat(ctx, req, func(resp api.ChatResponse) error {
		fullResponse = resp.Message.Content
		return nil
	})
	if err != nil {
		return "", err
	}

	return fullResponse, nil
}
