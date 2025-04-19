package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestStatusWindowUpdate(t *testing.T) {
	sizeMsg := tea.WindowSizeMsg{
		Width:  16,
		Height: 9,
	}

	ren := &lipgloss.Style{}
	w := StatusWindow{Renderer: ren, Focus: true}

	switch u, _ := w.Update(sizeMsg); u := u.(type) {
	case StatusWindow:
		assert.Equal(t, 16, u.Width)
		assert.Equal(t, 9, u.Height)

	default:
		t.Fatal("Expected the component to be a `win.StatusWindow`")
	}
}
