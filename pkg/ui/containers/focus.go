package containers

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
	win "github.com/vieitesss/egit/pkg/ui/windows"
)

func currentFocusTrace(trace []int, head comp.ComponentI) []int {
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
	return currentFocusTrace(trace, h.Components[h.CompFocused])
}

func (m *Container) setFocus(trace []int, focus bool) {
	if len(trace) == 0 {
		panic("updateFocus: trace len should be at least 1")
	}

	current := trace[0]
	if focus {
		m.CompFocused = current
	} else {
		m.CompFocused = -1
	}

	w, ok := m.Components[current].(*win.Window)
	if ok {
		if w.Focus == focus {
			panic("updateFocus: the window focus state is already set to the requested value")
		}
		w.Focus = focus
		m.Components[current] = w
		return
	}

	c, _ := m.Components[current].(*Container)
	c.setFocus(trace[1:], focus)
	m.Components[current] = c
}

func (m *Container) updateFocus(newTrace []int) {
	if len(newTrace) == 0 {
		panic("updateFocus: newTrace len should be at least 1")
	}

	currentTrace := currentFocusTrace([]int{}, m)

	if len(currentTrace) == 0 {
		panic("updateFocus: currentTrace len should be at least 1")
	}

	if reflect.DeepEqual(currentTrace, newTrace) {
		return
	}

	var i int
	c := m

	for currentTrace[i] == newTrace[i] {
		if x, ok := c.Components[currentTrace[i]].(*Container); ok {
			c = x
		} else {
			panic(fmt.Sprintf("updateFocus: the Container child should be a Container [new: %v, current: %v, ind: %d]", newTrace, currentTrace, i))
		}
		i++
	}

	c.setFocus(currentTrace[i:], false)
	c.setFocus(newTrace[i:], true)
}

func (m *Container) lastComponent(trace []int) comp.ComponentI {
	c := m

	for _, i := range trace {
		if child, ok := c.Components[i].(*Container); ok {
			c = child
		} else {
			return c.Components[i]
		}
	}

	return c
}

func (m *Container) findFirstWindow() []int {
	i := []int{0}

	c, ok := m.Components[0].(*Container)
	for ok {
		i = append(i, 0)
		c, ok = c.Components[0].(*Container)
	}

	return i
}

func (m *Container) isValidTrace(trace []int) bool {
	c := m

	for _, i := range trace {
		if i >= len(c.Components) {
			return false
		}

		if _, ok := c.Components[i].(*win.Window); ok {
			return true
		}

		c, _ = c.Components[i].(*Container)
	}

	return true
}

func (m *Container) changeToNextRow(trace []int) (bool, []int) {
	n := len(trace)

	if n < 2 {
		return false, []int{}
	}

	if _, ok := m.lastComponent(trace).(*Container); !ok {
		panic("changeToNextRow: the last component should be a Container")
	}

	if n%2 != 0 {
		panic("changeToNextRow: the last component should be a Column Container")
	}

	var newTrace []int
	newTrace = append(newTrace, trace[:n-1]...)
	newTrace[n-2]++
	valid := m.isValidTrace(newTrace)

	if !valid {
		return m.changeToNextRow(trace[:n-2])
	}

	return true, newTrace
}

func (m *Container) getNewTraceDown() []int {
	trace := currentFocusTrace([]int{}, m)
	n := len(trace)

	if n == 0 {
		panic("getNewTraceDown: trace len should be at least one")
	}

	if _, ok := m.lastComponent(trace).(*Container); ok {
		panic("getNewTraceDown: the trace should not end in a Container")
	}

	var (
		traceTilCont []int
		newTrace     []int
		ok           bool
	)

	traceTilCont = append(traceTilCont, trace[:n-1]...)
	c, _ := m.lastComponent(traceTilCont).(*Container)
	newFocus := trace[n-1] + 1

	if newFocus == len(c.Components) {
		ok, newTrace = m.changeToNextRow(traceTilCont)

		if !ok {
			return trace
		}

		var com comp.ComponentI
		com = m.lastComponent(newTrace)
		if _, ok = com.(*win.Window); ok {
			return newTrace
		}

		c, ok = com.(*Container)
		return append(newTrace, c.findFirstWindow()...)
	}

	last := c.Components[newFocus]

	if c, ok = last.(*Container); ok {
		newTrace = append(traceTilCont, newFocus)
		newTrace = append(newTrace, c.findFirstWindow()...)
		return newTrace
	}

	// The next window in the same container
	return append(traceTilCont, newFocus)
}

func (m *Container) changeFocus(msg tea.KeyMsg) {
	var newTrace []int

	switch msg.String() {
	case "left":
	case "right":
	case "down":
		newTrace = m.getNewTraceDown()
	case "up":

	default:
		panic("changeFocus: msg.String() should be one of 'left', 'right', 'down' or 'up'")
	}
	m.updateFocus(newTrace)
}
