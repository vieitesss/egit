package component

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ComponentI interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (ComponentI, tea.Cmd)
	View() string
	IsFocused() bool
	Size() (int, int)
	GetFixedHeight() int
}

type Component struct {
	Renderer *lipgloss.Style
	Width    int
	Height   int
	Focus    bool
}

func (c Component) IsFocused() bool {
	return c.Focus
}

func (c Component) Size() (int, int) {
	return c.Width, c.Height
}

func (c Component) GetFixedHeight() int {
	return 0
}
