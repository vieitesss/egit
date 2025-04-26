package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	"github.com/vieitesss/egit/pkg/ui/component"
	cont "github.com/vieitesss/egit/pkg/ui/containers"
	win "github.com/vieitesss/egit/pkg/ui/windows"
)
func TestTypes(t *testing.T) {
	ren := &lipgloss.Style{}
	layout := cont.NewColumn(
		ren,
		cont.NewRow(
			ren,
			cont.NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, true, 3),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
			cont.NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
		),
		win.NewWindow(ren, win.StatusWindow{}, false, 0),
	)

	assert.IsType(t, cont.Column{}, layout.Type)

	assert.Equal(t, 2, len(layout.Components))

	row := layout.Components[0].(cont.Container)
	assert.IsType(t, cont.Row{}, row.Type)
	w := layout.Components[1]
	assert.IsType(t, win.Window{}, w)

	assert.Equal(t, 2, len(row.Components))

	col1 := row.Components[0].(cont.Container)
	assert.IsType(t, cont.Column{}, col1.Type)
	col2 := row.Components[1].(cont.Container)
	assert.IsType(t, cont.Column{}, col2.Type)

	assert.Equal(t, 3, len(col1.Components))

	win1 := col1.Components[0].(win.Window)
	assert.IsType(t, win.StatusWindow{}, win1.Type)
	win2 := col1.Components[1].(win.Window)
	assert.IsType(t, win.StatusWindow{}, win2.Type)
	win3 := col1.Components[2].(win.Window)
	assert.IsType(t, win.StatusWindow{}, win3.Type)

	win4 := col2.Components[0].(win.Window)
	assert.IsType(t, win.StatusWindow{}, win4.Type)
	win5 := col2.Components[1].(win.Window)
	assert.IsType(t, win.StatusWindow{}, win5.Type)
}

func TestSizes(t *testing.T) {
	ren := &lipgloss.Style{}
	layout := cont.NewColumn(
		ren,
		cont.NewRow(
			ren,
			cont.NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, true, 3),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
			cont.NewColumn(
				ren,
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
				win.NewWindow(ren, win.StatusWindow{}, false, 0),
			),
		),
		win.NewWindow(ren, win.StatusWindow{}, false, 0),
	)

	var u component.ComponentI
	u, _ = layout.Update(tea.WindowSizeMsg{100, 100})
	if up, ok := u.(cont.Container); ok {
		layout = up
	}

	// Main
	main := layout.Component
	assert.Equal(t, 100, main.Height)
	assert.Equal(t, 100, main.Width)
	assert.Equal(t, 0, main.GetFixedHeight())

	// Row
	row := layout.Components[0].(cont.Container)
	assert.Equal(t, 50, row.Height)
	assert.Equal(t, 100, row.Width)
	assert.Equal(t, 0, row.GetFixedHeight())
	// Win
	win1 := layout.Components[1].(win.Window)
	assert.Equal(t, 50, win1.Height)
	assert.Equal(t, 100, win1.Width)
	assert.Equal(t, 0, win1.GetFixedHeight())

	// Col
	col1 := row.Components[0].(cont.Container)
	assert.Equal(t, 50, col1.Height)
	assert.Equal(t, 50, col1.Width)
	assert.Equal(t, 3, col1.GetFixedHeight())
	// Col
	col2 := row.Components[1].(cont.Container)
	assert.Equal(t, 50, col2.Height)
	assert.Equal(t, 50, col2.Width)
	assert.Equal(t, 0, col2.GetFixedHeight())

	// Win
	win2 := col1.Components[0].(win.Window)
	assert.Equal(t, 3, win2.Height)
	assert.Equal(t, 50, win2.Width)
	assert.Equal(t, 3, win2.GetFixedHeight())
	// Win
	win3 := col1.Components[1].(win.Window)
	assert.Equal(t, 24, win3.Height)
	assert.Equal(t, 50, win3.Width)
	assert.Equal(t, 0, win3.GetFixedHeight())
	// Win
	win4 := col1.Components[2].(win.Window)
	assert.Equal(t, 23, win4.Height)
	assert.Equal(t, 50, win4.Width)
	assert.Equal(t, 0, win4.GetFixedHeight())

	// Win
	win5 := col2.Components[0].(win.Window)
	assert.Equal(t, 25, win5.Height)
	assert.Equal(t, 50, win5.Width)
	assert.Equal(t, 0, win5.GetFixedHeight())
	// Win
	win6 := col2.Components[1].(win.Window)
	assert.Equal(t, 25, win6.Height)
	assert.Equal(t, 50, win6.Width)
	assert.Equal(t, 0, win6.GetFixedHeight())
}
