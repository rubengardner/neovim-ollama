package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Border   lipgloss.Style
	Input    lipgloss.Style
	Output   lipgloss.Style
	Prompt   lipgloss.Style
	Response lipgloss.Style
	Divider  string
	Selected lipgloss.Style
	File     lipgloss.Style
	Folder   lipgloss.Style
	Checked  lipgloss.Style
	Help     lipgloss.Style
}

func NewStyles() Styles {
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63"))

	return Styles{
		Border:   border,
		Input:    border.Padding(0, 1),
		Output:   border.Padding(0, 1).MarginBottom(1),
		Prompt:   lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Bold(true),
		Response: lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F")),
		Divider:  lipgloss.NewStyle().Foreground(lipgloss.Color("#444")).Render(strings.Repeat("─", 40)),
		Selected: lipgloss.NewStyle().Background(lipgloss.Color("#3C3836")).Foreground(lipgloss.Color("#EBDBB2")),
		File:     lipgloss.NewStyle().Foreground(lipgloss.Color("#83A598")),
		Folder:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FABD2F")).Bold(true),
		Checked:  lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).Bold(true),
		Help:     lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")).Italic(true),
	}
}
