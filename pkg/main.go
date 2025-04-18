package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/egit/pkg/ui"
)

func main() {
	pro := tea.NewProgram(ui.NewEgit(), tea.WithAltScreen())

	if _, err := pro.Run(); err != nil {
		fmt.Errorf("Error initializing Egit: %v", err)
	}
}
