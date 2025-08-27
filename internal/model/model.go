package model

import (
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rubengardner/neovim-ollama/internal/files"
)

type Mode int

const (
	ChatMode Mode = iota
	FileSelectMode
	ReviewMode
)

type Model struct {
	Input           textinput.Model
	Viewport        viewport.Model
	Width           int
	Height          int
	IsWaiting       bool
	Err             error
	Spinner         spinner.Model
	History         []files.ChatMessage
	Mode            Mode
	Files           []files.FileItem
	FilesCursor     int
	CurrentDir      string
	SelectedFiles   []string
	FilesViewport   viewport.Model
	ProposedChanges []files.FileChange
	ReviewCursor    int
	ReviewViewport  viewport.Model
}

type (
	ResponseMsg      string
	ErrorMsg         error
	ReviewChangesMsg []files.FileChange
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
		ProposedChanges: []files.FileChange{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.Spinner.Tick)
}
