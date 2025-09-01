package files

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Files         []FileItem
	Cursor        int
	CurrentDir    string
	SelectedFiles []string
	Viewport      viewport.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() Model {
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	return Model{
		Files:         []FileItem{},
		CurrentDir:    currentDir,
		SelectedFiles: []string{},
		Viewport:      viewport.New(80, 20),
	}
}

type FileItem struct {
	Name     string
	Path     string
	IsDir    bool
	Selected bool
}

type FileChange struct {
	FilePath     string
	OriginalCode string
	ProposedCode string
	Description  string
	Accepted     *bool
	LineStart    int
	LineEnd      int
}
