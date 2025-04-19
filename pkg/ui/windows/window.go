package windows

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/charmbracelet/lipgloss/v2"
)

type WindowI interface {
	handleKeyPress(tea.KeyPressMsg) tea.Cmd
	updateWindowSize(w, h int)
}

type Window struct {
	viewport viewport.Model
	Renderer *lipgloss.Style
	Width    int
	Height   int
	Focus    bool
	Content  string
}

func (m *Window) updateWindowSize(w, h int) {
	m.Width = w
	m.Height = h

	m.viewport.SetWidth(m.Width - 2)
	m.viewport.SetHeight(m.Height - 2)
}
