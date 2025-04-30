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
			NewColumn(
				ren,
				NewRow(
					ren,
					NewColumn(
						ren,
						win.NewWindow(ren, win.LogWindow{}, true, 0),
					),
				),
			),
		),
		win.NewWindow(ren, win.LogWindow{}, false, 0),
	)

	trace = focusTrace([]int{}, l)
	assert.Equal(t, []int{0, 0, 0, 0, 0}, trace)

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

func TestUpdateFocus(t *testing.T) {
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

func TestNewTraceDown(t *testing.T) {
	ren := &lipgloss.Style{}
	l := NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 3),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
				win.NewWindow(ren, win.LogWindow{}, true, 0),
			),
			NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
		),
	)

	newTrace := l.getNewTraceDown()
	assert.Equal(t, []int{0, 0, 2}, newTrace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.LogWindow{}, true, 0),
			),
		),
		NewRow(
			ren,
			NewColumn(
				ren,
				NewRow(
					ren,
					NewColumn(
						ren,
						win.NewWindow(ren, win.LogWindow{}, false, 0),
					),
				),
			),
		),
	)

	newTrace = l.getNewTraceDown()
	assert.Equal(t, []int{1, 0, 0, 0, 0}, newTrace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.LogWindow{}, true, 0),
				win.NewWindow(ren, win.LogWindow{}, false, 0),
			),
		),
	)

	newTrace = l.getNewTraceDown()
	assert.Equal(t, []int{0, 0, 1}, newTrace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				NewRow(
					ren,
					NewColumn(
						ren,
						win.NewWindow(ren, win.LogWindow{}, true, 0),
					),
				),
			),
		),
		win.NewWindow(ren, win.LogWindow{}, false, 0),
	)

	newTrace = l.getNewTraceDown()
	assert.Equal(t, []int{1}, newTrace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.LogWindow{}, true, 0),
			),
		),
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.LogWindow{}, false, 0),
			),
		),
	)

	newTrace = l.getNewTraceDown()
	assert.Equal(t, []int{1, 0, 0}, newTrace)

	l = NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, true, 0),
				NewRow(
					ren,
					NewColumn(
						ren,
						win.NewWindow(ren, win.StatusWindow{}, false, 0),
						win.NewWindow(ren, win.LogWindow{}, false, 0),
					),
				),
			),
		),
	)

	newTrace = l.getNewTraceDown()
	assert.Equal(t, []int{0, 0, 1, 0, 0}, newTrace)
}

func TestLastComponent(t *testing.T) {
	ren := &lipgloss.Style{}
	l := NewColumn(
		ren,
		NewRow(
			ren,
			NewColumn(
				ren,
				win.NewWindow(ren, win.BranchWindow{}, false, 3),
				win.NewWindow(ren, win.LogWindow{}, true, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
			NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
		),
	)

	c := l.lastComponent([]int{0, 0, 1})
	assert.IsType(t, &win.Window{}, c)

	c = l.lastComponent([]int{0, 0})
	assert.IsType(t, &Container{}, c)
	col, _ := c.(*Container)
	assert.IsType(t, Column{}, col.Type)
}
