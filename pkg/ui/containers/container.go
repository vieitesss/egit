package containers

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
)

type Container interface {
	comp.Component
	getViewComponents() string
	getWindowSizeDimensions(cmps, w, h int) (width, height int)
}

type baseContainer struct {
	self       Container
	Renderer   *lipgloss.Style
	Components []comp.Component
	Width      int
	Height     int
	focused    int
}

type Column struct {
	baseContainer
}

type Row struct {
	baseContainer
}

func NewColumn(ren *lipgloss.Style, c ...comp.Component) Column {
	col := Column{newBaseContainer(ren, c...)}
	col.self = &col
	return col
}

func NewRow(ren *lipgloss.Style, c ...comp.Component) Row {
	row := Row{newBaseContainer(ren, c...)}
	row.self = &row
	return row
}

func newBaseContainer(ren *lipgloss.Style, c ...comp.Component) baseContainer {
	var cmps []comp.Component
	focus := -1

	for i, cmp := range c {
		if ok := cmp.IsFocused(); ok {
			focus = i
		}
		cmps = append(cmps, cmp)
	}

	return baseContainer{
		Renderer:   ren,
		Components: cmps,
		focused:    focus,
	}
}

func (m *baseContainer) updateAll(update bool, msg tea.Msg) tea.Cmd {
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

func (m baseContainer) Init() tea.Cmd {
	return m.updateAll(false, nil)
}

func (m baseContainer) Update(msg tea.Msg) (comp.Component, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd tea.Cmd
	)


	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		cmps := len(m.Components)
		for i, c := range m.Components {
			width, height := m.self.getWindowSizeDimensions(cmps, m.Width, m.Height)
			m.Components[i], cmd = c.Update(tea.WindowSizeMsg{
				Width:  width,
				Height: height,
			})
			cmds = append(cmds, cmd)
		}

	case tea.KeyPressMsg:
		m.Components[m.focused], cmd = m.Components[m.focused].Update(msg)
		cmds = append(cmds, cmd)

	default:
		return m, m.updateAll(true, msg)
	}

	return m, tea.Batch(cmds...)
}

func (m baseContainer) View() string {
	return m.Renderer.
		Width(m.Width).
		Height(m.Height).
		Render(m.self.getViewComponents())
}

func (m baseContainer) IsFocused() bool {
	return m.focused != -1
}

func (m Column) getViewComponents() string {
	var views []string
	for _, c := range m.Components {
		views = append(views, c.View())
	}
	return lipgloss.JoinVertical(0, views...)
}

func (m Row) getViewComponents() string {
	var views []string
	for _, c := range m.Components {
		views = append(views, c.View())
	}
	return lipgloss.JoinHorizontal(0, views...)
}

func (m Column) getWindowSizeDimensions(cmps, w, h int) (width, height int) {
	return w, h / cmps
}

func (m Row) getWindowSizeDimensions(cmps, w, h int) (width, height int) {
	return w / cmps, h
}
