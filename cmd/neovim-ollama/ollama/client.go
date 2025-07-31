package ollama

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func Generate(prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var genResp GenerateResponse
	if err := json.Unmarshal(bodyBytes, &genResp); err != nil {
		return "", err
	}

	return genResp.Response, nil
}
