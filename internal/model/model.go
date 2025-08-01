package model

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
}

type (
	ResponseMsg      string
	ErrorMsg         error
	FilesLoadedMsg   []FileItem
	ReviewChangesMsg []FileChange
)
