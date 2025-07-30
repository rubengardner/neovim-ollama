package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type StreamChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type StreamCallback func(chunk StreamChunk)

func StreamGenerate(prompt string, onChunk StreamCallback) error {
	reqBody := GenerateRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: true,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	for {
		var chunk StreamChunk
		if err := decoder.Decode(&chunk); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		onChunk(chunk)
		if chunk.Done {
			break
		}
	}

	return nil
}
