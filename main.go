package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/mrSamDev/brew-potato/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
