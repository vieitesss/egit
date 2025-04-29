package containers

import (
	"testing"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	win "github.com/vieitesss/egit/pkg/ui/windows"
)

func TestTrace(t *testing.T) {
	ren := &lipgloss.Style{}
	l := NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 3),
				win.NewWindow(ren, win.StatusWindow{}, true, 0),
			),
		),
	)

	trace := focusTrace([]int{}, l)
	assert.Equal(t, []int{0, 0, 1}, trace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 0),
			),
			NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, true, 0),
			),
		),
	)

	trace = focusTrace([]int{}, l)
	assert.Equal(t, []int{0, 1, 0}, trace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			win.NewWindow(ren, win.BranchWindow{}, false, 3),
		),
		win.NewWindow(ren, win.StatusWindow{}, true, 3),
	)

	trace = focusTrace([]int{}, l)
	assert.Equal(t, []int{1}, trace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 3),
			),
			NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
		),
		win.NewWindow(ren, win.StatusWindow{}, false, 3),
	)

	assert.Panics(t, func() { focusTrace([]int{}, l) })
}

func TestSetFocus(t *testing.T) {
	ren := &lipgloss.Style{}
	l := NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
		),
		win.NewWindow(ren, win.StatusWindow{}, false, 0),
	)

	l.setFocus([]int{0, 0, 1})
	trace := focusTrace([]int{}, l)
	assert.Equal(t, []int{0, 0, 1}, trace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
		),
		win.NewWindow(ren, win.StatusWindow{}, false, 0),
	)

	l.setFocus([]int{1})
	trace = focusTrace([]int{}, l)
	assert.Equal(t, []int{1}, trace)
}

func TestGetContainer(t *testing.T) {
	ren := &lipgloss.Style{}
	l := NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 3),
				win.NewWindow(ren, win.StatusWindow{}, true, 0),
			),
			win.NewWindow(ren, win.BranchWindow{}, false, 3),
		),
	)

	trace := []int{0, 0, 0}
	elem := l.getLastContainer(trace)
	assert.IsType(t, Column{}, elem.Type)
	assert.Equal(t, 2, len(elem.Components))

	trace = []int{0, 1}
	elem = l.getLastContainer(trace)
	assert.IsType(t, Row{}, elem.Type)
	assert.Equal(t, 2, len(elem.Components))

	trace = []int{0}
	elem = l.getLastContainer(trace)
	assert.IsType(t, Column{}, elem.Type)
	assert.Equal(t, *l, elem)
	assert.Equal(t, 1, len(elem.Components))

	trace = []int{}
	assert.Panics(t, func() { l.getLastContainer(trace) })

	trace = []int{0, 2}
	assert.Panics(t, func() { l.getLastContainer(trace) })
}

func TestRemoveFocus(t *testing.T) {
	ren := &lipgloss.Style{}
	l := NewColumn(
		ren,
		win.NewWindow(ren, win.BranchWindow{}, false, 3),
		win.NewWindow(ren, win.StatusWindow{}, true, 0),
	)

	l.removeFocus([]int{1})
	assert.Equal(t, -1, l.CompFocused)
	w, _ := l.Components[1].(*win.Window)
	assert.False(t, w.Focus)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 3),
				win.NewWindow(ren, win.StatusWindow{}, true, 0),
			),
			win.NewWindow(ren, win.BranchWindow{}, false, 3),
		),
	)

	l.removeFocus([]int{0, 0, 1})
	assert.Equal(t, -1, l.CompFocused)
	c1, _ := l.Components[0].(*Container)
	assert.Equal(t, -1, c1.CompFocused)
	c2, _ := c1.Components[0].(*Container)
	assert.Equal(t, -1, c2.CompFocused)
	w, _ = c2.Components[1].(*win.Window)
	assert.False(t, w.Focus)

	assert.Panics(t, func() { l.removeFocus([]int{}) })
	assert.Panics(t, func() { l.removeFocus([]int{0, 1}) })
}

func TestUpdateFocus(t *testing.T){
	ren := &lipgloss.Style{}
	l := NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 3),
				win.NewWindow(ren, win.StatusWindow{}, true, 0),
			),
			win.NewWindow(ren, win.BranchWindow{}, false, 3),
		),
	)

	l.updateFocus([]int{0, 1})
	trace := focusTrace([]int{}, l)
	assert.Equal(t, []int{0, 1}, trace)
}
