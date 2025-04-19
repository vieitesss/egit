package containers

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
)

type Container struct {
	base      baseContainer
	joiner    func([]string) string
	dimension func(count, w, h int) (int, int)
}

type baseContainer struct {
	Renderer   *lipgloss.Style
	Components []comp.Component
	Width      int
	Height     int
	focused    int
}

func newBaseContainer(ren *lipgloss.Style, cmps ...comp.Component) baseContainer {
	b := baseContainer{
		Renderer:   ren,
		Components: make([]comp.Component, len(cmps)),
		focused:    -1,
	}
	for i, c := range cmps {
		b.Components[i] = c
		if c.IsFocused() {
			b.focused = i
		}
	}
	return b
}

func (m baseContainer) updateAll(doUpdate bool, msg tea.Msg) (baseContainer, tea.Cmd) {
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

func NewColumn(ren *lipgloss.Style, cmps ...comp.Component) Container {
	return Container{
		base:      newBaseContainer(ren, cmps...),
		joiner:    func(vs []string) string { return lipgloss.JoinVertical(0, vs...) },
		dimension: func(n, w, h int) (int, int) { return w, h / n },
	}
}

func NewRow(ren *lipgloss.Style, cmps ...comp.Component) Container {
	return Container{
		base:      newBaseContainer(ren, cmps...),
		joiner:    func(vs []string) string { return lipgloss.JoinHorizontal(0, vs...) },
		dimension: func(n, w, h int) (int, int) { return w / n, h },
	}
}

func (m Container) Init() tea.Cmd {
	var cmd tea.Cmd
	m.base, cmd = m.base.updateAll(false, nil)
	return cmd
}

func (m Container) Update(msg tea.Msg) (comp.Component, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.base.Width, m.base.Height = msg.Width, msg.Height
		n := len(m.base.Components)
		for i, child := range m.base.Components {
			w, h := m.dimension(n, msg.Width, msg.Height)
			m.base.Components[i], cmd = child.Update(tea.WindowSizeMsg{Width: w, Height: h})
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		if m.base.focused >= 0 && m.base.focused < len(m.base.Components) {
			m.base.Components[m.base.focused], cmd = m.base.Components[m.base.focused].Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		m.base, cmd = m.base.updateAll(true, msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Container) View() string {
	return m.base.
		Renderer.
		Width(m.base.Width).
		Height(m.base.Height).
		Render(m.getViewComponents())
}

func (m Container) Size() (int, int) {
	return m.base.Width, m.base.Height
}

func (m Container) IsFocused() bool {
	return m.base.focused != -1
}

func (m Container) getViewComponents() string {
	var vs []string

	for _, child := range m.base.Components {
		vs = append(vs, child.View())
	}

	return m.joiner(vs)
}
