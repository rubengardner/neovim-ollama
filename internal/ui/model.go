package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
)

type Model struct {
	Width     int
	Height    int
	IsWaiting bool
	Err       error
	Spinner   spinner.Model
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		Width:     80,
		Height:    24,
		IsWaiting: false,
		Spinner:   s,
	}
}
