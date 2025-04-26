package containers

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
	win "github.com/vieitesss/egit/pkg/ui/windows"
)

type ContainerType interface {
	Joiner([]string) string
	Dimension(count, w, h int) (int, int)
}

type Container struct {
	comp.Component
	Type        ContainerType
	Components  []comp.ComponentI
	CompFocused int
}

func newContainer(ren *lipgloss.Style, cmps ...comp.ComponentI) Container {
	con := Container{
		Component:   comp.Component{Renderer: ren},
		Components:  make([]comp.ComponentI, len(cmps)),
		CompFocused: -1,
	}

	for i, c := range cmps {
		con.Components[i] = c
		if c.IsFocused() {
			con.CompFocused = i
		}
	}

	return con
}

func (m Container) updateAll(doUpdate bool, msg tea.Msg) (Container, tea.Cmd) {
	var cmds []tea.Cmd

	if doUpdate {
		for i, c := range m.Components {
			updated, cmd := c.Update(msg)
			m.Components[i] = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	} else {
		for _, c := range m.Components {
			cmds = append(cmds, c.Init())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Container) Init() tea.Cmd {
	var cmd tea.Cmd
	m, cmd = m.updateAll(false, nil)
	return cmd
}

func (m *Container) updateSize(w, h int) tea.Cmd {
	m.Width, m.Height = w, h

	tiledHeight, elemsToTile := m.tileParams()
	cmds := m.updateChildrenSize(tiledHeight, elemsToTile)

	return cmds
}

func (m Container) tileParams() (int, int) {
	tiledHeight := m.Height - m.GetFixedHeight()
	elemsToTile := len(m.Components)
	for _, e := range m.Components {
		if w, ok := e.(win.Window); ok && w.MaxHeight > 0 {
			elemsToTile--
		}
	}

	return tiledHeight, elemsToTile
}

func (m *Container) updateChildrenSize(tiledHeight, elemsToTile int) tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	remainSpace := tiledHeight % elemsToTile

	for i, child := range m.Components {
		width, height := m.Type.Dimension(elemsToTile, m.Width, tiledHeight)
		if c, ok := child.(win.Window); ok && c.MaxHeight == 0 && remainSpace > 0 {
			height += 1
			remainSpace -= 1
		}
		m.Components[i], cmd = child.Update(tea.WindowSizeMsg{Width: width, Height: height})
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m Container) Update(msg tea.Msg) (comp.ComponentI, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmds = append(cmds, m.updateSize(msg.Width, msg.Height))

	case tea.KeyMsg:
		if m.CompFocused >= 0 && m.CompFocused < len(m.Components) {
			m.Components[m.CompFocused], cmd = m.Components[m.CompFocused].Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		m, cmd = m.updateAll(true, msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Container) View() string {
	return m.
		Renderer.
		Width(m.Width).
		Height(m.Height).
		Render(m.getViewComponents())
}

func (m Container) Size() (int, int) {
	return m.Width, m.Height
}

func (m Container) IsFocused() bool {
	return m.CompFocused != -1
}

func (m Container) getViewComponents() string {
	var vs []string

	for _, child := range m.Components {
		vs = append(vs, child.View())
	}

	return m.Type.Joiner(vs)
}

func (m Container) GetFixedHeight() int {
	var fixed int
	for _, c := range m.Components {
		if c, ok := c.(win.Window); ok {
			fixed += c.MaxHeight
		}
	}

	return fixed
}
