package component

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Component interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Component, tea.Cmd)
	View() string
	IsFocused() bool
	Size() (int, int)
}
