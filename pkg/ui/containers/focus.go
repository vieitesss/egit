package containers

import (
	"reflect"

	tea "github.com/charmbracelet/bubbletea/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
	win "github.com/vieitesss/egit/pkg/ui/windows"
)

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
