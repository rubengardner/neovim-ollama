package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
