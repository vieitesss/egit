package containers

import (
	"reflect"

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

func newContainer(ren *lipgloss.Style, cmps ...comp.ComponentI) *Container {
	con := &Container{
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

func (m *Container) updateAll(doUpdate bool, msg tea.Msg) tea.Cmd {
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

	return tea.Batch(cmds...)
}

func (m *Container) Init() tea.Cmd {
	var cmd tea.Cmd
	cmd = m.updateAll(false, nil)
	return cmd
}

func (m Container) tileParams() (int, int) {
	tiledHeight := m.Height - m.GetFixedHeight()
	elemsToTile := len(m.Components)
	for _, e := range m.Components {
		if w, ok := e.(*win.Window); ok && w.MaxHeight > 0 {
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
		if c, ok := child.(*win.Window); ok && c.MaxHeight == 0 && remainSpace > 0 {
			height += 1
			remainSpace -= 1
		}
		m.Components[i], cmd = child.Update(tea.WindowSizeMsg{Width: width, Height: height})
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *Container) updateSize(w, h int) tea.Cmd {
	m.Width, m.Height = w, h

	tiledHeight, elemsToTile := m.tileParams()
	cmds := m.updateChildrenSize(tiledHeight, elemsToTile)

	return cmds
}

func focusTrace(trace []int, head comp.ComponentI) []int {
	var (
		h  *Container
		ok bool
	)

	if h, ok = head.(*Container); !ok {
		return trace
	}

	if h.CompFocused == -1 {
		panic("focusTrace: there is nothing focused")
	}

	trace = append(trace, h.CompFocused)
	return focusTrace(trace, h.Components[h.CompFocused])
}

func (m *Container) getLastContainer(trace []int) Container {
	n := len(trace)
	if n == 0 {
		panic("getContainer: the trace does not have elements")
	}

	if n < 2 {
		return *m
	}

	var i int
	head := m

	h, ok := head.Components[trace[i]].(*Container)
	for ok {
		head = h
		i++
		h, ok = head.Components[trace[i]].(*Container)
	}

	return *head
}

func (m *Container) removeFocus(trace []int) {
	if len(trace) == 0 {
		panic("removeFocus: trace len should be at least 1")
	}

	m.CompFocused = -1
	current := trace[0]

	w, ok := m.Components[current].(*win.Window)
	if ok {
		if !w.Focus {
			panic("removeFocus: the window should be focused")
		}
		w.Focus = false
		m.Components[current] = w
		return
	}

	c, _ := m.Components[current].(*Container)
	c.removeFocus(trace[1:])
	m.Components[current] = c
	return
}

func (m *Container) setFocus(trace []int) {
	if len(trace) == 0 {
		panic("setFocus: newTrace len should be at least 1")
	}

	current := trace[0]
	m.CompFocused = current

	w, ok := m.Components[current].(*win.Window)
	if ok {
		if w.Focus {
			panic("setFocus: the window should not be focused")
		}
		w.Focus = true
		m.Components[current] = w
		return
	}

	c, _ := m.Components[current].(*Container)
	c.setFocus(trace[1:])
	m.Components[current] = c
	return
}

func (m *Container) updateFocus(newTrace []int) {
	if len(newTrace) == 0 {
		panic("updateFocus: newTrace len should be at least 1")
	}

	currentTrace := focusTrace([]int{}, m)

	if len(currentTrace) == 0 {
		panic("updateFocus: currentTrace len should be at least 1")
	}

	if reflect.DeepEqual(currentTrace, newTrace) {
		panic("updateFocus: newTrace should not be equal to currentTrace")
	}

	var i int
	c := m

	// new := []int{0, 1}
	// cur := []int{0, 0, 1}
	for currentTrace[i] == newTrace[i] {
		if x, ok := c.Components[currentTrace[i]].(*Container); ok {
			c = x
		} else {
			panic("updateFocus: the Container child should be a Container")
		}
		i++
	}

	c.removeFocus(currentTrace[i:])

	if c.CompFocused != -1 {
		panic("updateFocus: removeFocus did not work")
	}

	c.setFocus(newTrace[i:])

	currentTrace = focusTrace([]int{}, m)
	if !reflect.DeepEqual(currentTrace, newTrace) {
		panic("updateFocus: update did not work")
	}
}

func (m *Container) changeFocus(msg tea.KeyMsg) {
	trace := focusTrace([]int{}, m)
	n := len(trace)

	switch msg.String() {
	case "left":
		// if n > 1 && trace[n-2] > 0 {
		// 	trace = append(trace[:n-2], trace[n-2]-1)
		// 	setFocus(trace, m)
		// }
	case "right":
	case "down":
		c := m.getLastContainer(trace)
		if len(c.Components) > trace[n-1]+1 {
			newTrace := append(trace[:n-1], trace[n-1]+1)
			m.updateFocus(newTrace)
		}
	case "up":

	default:
		panic("changeFocus: msg.String() should be one of 'left', 'right', 'down' or 'up'")
	}
}

func (m *Container) Update(msg tea.Msg) (comp.ComponentI, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmds = append(cmds, m.updateSize(msg.Width, msg.Height))

	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "down", "up":
			m.changeFocus(msg)
			return m, nil
		}
		if m.CompFocused >= 0 && m.CompFocused < len(m.Components) {
			m.Components[m.CompFocused], cmd = m.Components[m.CompFocused].Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		cmd = m.updateAll(true, msg)
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
		if c, ok := c.(*win.Window); ok {
			fixed += c.MaxHeight
		}
	}

	return fixed
}
