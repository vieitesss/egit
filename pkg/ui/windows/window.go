package windows

import (
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vieitesss/egit/pkg/cmd"
	"github.com/vieitesss/egit/pkg/ui/msgs"
)

type WindowI interface {
	handleKeyPress(tea.KeyPressMsg) tea.Cmd
	updateWindowSize(w, h int)
}

type Window struct {
	viewport  viewport.Model
	Renderer  *lipgloss.Style
	Width     int
	Height    int
	MaxHeight int
	Focus     bool
	Loaded    bool
	Content   string
}

func (m *Window) updateWindowSize(w, h int) {
	m.Width = w
	m.Height = h

	m.viewport.SetWidth(m.Width - 2)
	m.viewport.SetHeight(m.Height - 2)
}

func (m *Window) runGitCmdOut(c ...string) tea.Cmd {
	return func() tea.Msg {
		out, err := cmd.GitCmdOut(c...)
		return msgs.GitCmdOutMsg{
			Out: out,
			Err: err,
		}
	}
}
