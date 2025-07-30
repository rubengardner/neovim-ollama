package ui

import (
	"github.com/charmbracelet/glamour"
)

func renderMarkdown(input string) string {
	out, err := glamour.Render(input, "dark")
	if err != nil {
		return input
	}
	return out
}
