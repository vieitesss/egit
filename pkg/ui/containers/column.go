package container

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
)

type Container comp.Component

type Column struct {
	Container
	Renderer   *lipgloss.Style
	Components []comp.Component
	Width      int
	Height     int

	// The component that is focused in the column
	focused int
}

func NewColumn(ren *lipgloss.Style, c ...comp.Component) Column {
	var cmps []comp.Component
	focus := -1

	for i, cmp := range c {
		if ok := cmp.IsFocused(); ok {
			focus = i
		}
		cmps = append(cmps, cmp)
	}

	return Column{
		Renderer:   ren,
		Components: cmps,
		focused:    focus,
	}
}

func (m Column) updateAll(update bool, msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	if update {
		for i, c := range m.Components {
			m.Components[i], cmd = c.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		for _, c := range m.Components {
			cmd = c.Init()
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (m Column) Init() tea.Cmd {
	return m.updateAll(false, nil)
}

func (m Column) Update(msg tea.Msg) (comp.Component, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		cmps := len(m.Components)
		for i, c := range m.Components {
			m.Components[i], _ = c.Update(tea.WindowSizeMsg{
				Width:  m.Width,
				Height: m.Height / cmps,
			})
		}

	case tea.KeyPressMsg:
		m.Components[m.focused], cmd = m.Components[m.focused].Update(msg)

	default:
		m.updateAll(true, msg)
	}

	return m, cmd
}

func (m Column) View() string {
	var views []string

	for _, c := range m.Components {
		views = append(views, c.View())
	}

	return m.Renderer.
		Width(m.Width).
		Height(m.Height).
		Render(lipgloss.JoinVertical(0, views...))
}

func (m Column) IsFocused() bool {
	if m.focused == -1 {
		return false
	}

	return true
}
