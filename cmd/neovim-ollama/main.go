package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.Layout)

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, ui.SendPrompt); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		log.Panicln(err)
	}

	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
