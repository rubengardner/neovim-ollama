package ui

import (
	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(input string) string {
	out, err := glamour.Render(input, "tokyo-night")
	if err != nil {
		return input
	}
	return out
}
