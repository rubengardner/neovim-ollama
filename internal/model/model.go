package model

import (
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
)

type Mode int

const (
	ChatMode Mode = iota
	FileSelectMode
	ReviewMode
)

type ChatMessage struct {
	Role    string
	Content string
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
	Accepted     *bool // nil = pending, true = accepted, false = rejected
	LineStart    int
	LineEnd      int
}

type Model struct {
	Input           textinput.Model
	Viewport        viewport.Model
	Width           int
	Height          int
	IsWaiting       bool
	Err             error
	Spinner         spinner.Model
	History         []ChatMessage
	Mode            Mode
	Files           []FileItem
	FilesCursor     int
	CurrentDir      string
	SelectedFiles   []string
	FilesViewport   viewport.Model
	ProposedChanges []FileChange
	ReviewCursor    int
	ReviewViewport  viewport.Model
	Styles          ui.Styles
}

type (
	ResponseMsg      string
	ErrorMsg         error
	FilesLoadedMsg   []FileItem
	ReviewChangesMsg []FileChange
)

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter prompt (Ctrl+F for files, Ctrl+C to exit)"
	ti.Focus()
	ti.CharLimit = 500

	vp := viewport.New(0, 0)
	filesVp := viewport.New(0, 0)
	reviewVp := viewport.New(0, 0)

	sp := spinner.New()
	sp.Spinner = spinner.Line
	sp.Style = lipgloss.NewStyle()

	currentDir, _ := os.Getwd()

	return Model{
		Input:           ti,
		Viewport:        vp,
		FilesViewport:   filesVp,
		ReviewViewport:  reviewVp,
		Spinner:         sp,
		Mode:            ChatMode,
		CurrentDir:      currentDir,
		SelectedFiles:   []string{},
		ProposedChanges: []FileChange{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.Spinner.Tick)
}
