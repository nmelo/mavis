package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nmelo/typer/internal/app"
	"github.com/nmelo/typer/internal/drill"
	"github.com/nmelo/typer/internal/progress"
)

func main() {
	path := progress.DefaultPath()
	prog, err := progress.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading progress: %v\n", err)
		os.Exit(1)
	}

	words := drill.LoadWordList()
	model := app.New(prog, path, words)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
