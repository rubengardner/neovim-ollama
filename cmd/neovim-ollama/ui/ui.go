package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ollama"
)

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("output", 0, 0, maxX-1, maxY-5); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Wrap = true
		v.Autoscroll = true
		v.Title = "Model Response"
	}

	if v, err := g.SetView("input", 0, maxY-4, maxX-1, maxY-1); err != nil && err != gocui.ErrUnknownView {
		return err
	} else {
		v.Editable = true
		v.Title = "Enter Prompt"
		_, _ = g.SetCurrentView("input")
	}

	return nil
}

func SendPrompt(g *gocui.Gui, v *gocui.View) error {
	text := strings.TrimSpace(v.Buffer())
	v.Clear()
	v.SetCursor(0, 0)

	if text == "" {
		return nil
	}

	go func() {
		var fullResponse strings.Builder

		err := ollama.StreamGenerate(text, func(chunk ollama.StreamChunk) {
			fullResponse.WriteString(chunk.Response)
		})
		if err != nil {
			PrintToOutput(g, fmt.Sprintf("Error: %v\n", err))
			return
		}

		markdown := renderMarkdown(fullResponse.String())
		PrintToOutput(g, markdown)
	}()

	return nil
}

func PrintToOutput(g *gocui.Gui, text string) {
	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("output")
		v.Clear()
		fmt.Fprint(v, text)
		return nil
	})
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
