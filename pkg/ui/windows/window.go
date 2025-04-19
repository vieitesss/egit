package windows

import (
	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/charmbracelet/lipgloss/v2"
)

type Window struct {
	viewport viewport.Model
	Renderer *lipgloss.Style
	Width    int
	Height   int
	Focus    bool
	Content  string
}
