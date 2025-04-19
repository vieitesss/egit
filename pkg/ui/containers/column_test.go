package containers

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
	"github.com/vieitesss/egit/pkg/cmd"
	"github.com/vieitesss/egit/pkg/ui/msgs"
	win "github.com/vieitesss/egit/pkg/ui/windows"
)

func TestColumnCreation(t *testing.T) {
	ren := &lipgloss.Style{}
	f := win.Window{Renderer: ren, Focus: true}
	w := win.Window{Renderer: ren}
	sww := win.StatusWindow{w}
	swf := win.StatusWindow{f}

	col := NewColumn(ren, swf)

	assert.IsType(t, Container{}, col)
	assert.Equal(t, len(col.base.Components), 1)
	assert.True(t, col.IsFocused(), 0)

	col = NewColumn(ren, sww, swf)

	assert.Equal(t, len(col.base.Components), 2)
	assert.True(t, col.IsFocused(), 1)
}

func TestColumnInit(t *testing.T) {
	expect, err := cmd.GitCmdOut("status")

	ren := &lipgloss.Style{}
	w := win.StatusWindow{win.Window{Renderer: ren, Focus: true}}

	col := NewColumn(ren, w)
	cmds := col.Init()
	msg := cmds()

	switch msg := msg.(type) {
	case msgs.StatusCmdMsg:
		assert.Equal(t, expect, msg.Out)
		assert.Equal(t, err, msg.Err)

	default:
		t.Fatal("Expected the message to be a `msgs.StatusCmdMsg`")
	}
}

func TestColumnUpdate(t *testing.T) {
	statusMsg := msgs.StatusCmdMsg{
		Out: "hello",
		Err: error(nil),
	}

	ren := &lipgloss.Style{}
	w := win.StatusWindow{win.Window{Renderer: ren, Focus: true}}
	col := NewColumn(ren, w)

	switch u, _ := col.Update(statusMsg); u := u.(type) {
	case Container:
		switch x := u.base.Components[0]; x := x.(type) {
		case win.StatusWindow:
			assert.Equal(t, x.Content, "hello")

		default:
			t.Fatal("Expected the component to be a `win.StatusWindow`")
		}

	default:
		t.Fatal("Expected the component to be a `Column`")
	}

	sizeMsg := tea.WindowSizeMsg{
		Width:  16,
		Height: 9,
	}

	switch u, _ := col.Update(sizeMsg); u := u.(type) {
	case Container:
		assert.Equal(t, 16, u.base.Width)
		assert.Equal(t, 9, u.base.Height)
		switch x := u.base.Components[0]; x := x.(type) {
		case win.StatusWindow:
			assert.Equal(t, 16, x.Width)
			assert.Equal(t, 9, x.Height)

		default:
			t.Fatal("Expected the component to be a `win.StatusWindow`")
		}

	default:
		t.Fatal("Expected the component to be a `Column`")
	}
}
