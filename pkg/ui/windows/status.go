package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vieitesss/egit/pkg/cmd"
	comp "github.com/vieitesss/egit/pkg/ui/component"
	"github.com/vieitesss/egit/pkg/ui/msgs"
)

type Window comp.Component

type StatusWindow struct {
	Window
	viewport viewport.Model
	Renderer *lipgloss.Style
	Width    int
	Height   int
	Focus    bool
	Content  string
}

func (m StatusWindow) Init() tea.Cmd {
	m.viewport = viewport.New()

	return func() tea.Msg {
		out, err := cmd.GitCmdOut("status")
		return msgs.StatusCmdMsg{
			Out: out,
			Err: err,
		}
	}
}

func (m StatusWindow) Update(msg tea.Msg) (comp.Component, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.viewport.SetWidth(m.Width - 2)
		m.viewport.SetHeight(m.Height - 2)

	case tea.KeyPressMsg:
		m, cmd = m.handleKeyPress(msg)
		cmds = append(cmds, cmd)

	case msgs.StatusCmdMsg:
		if msg.Err != nil {
			m.Content = fmt.Sprintf("Error executing command: %v", msg.Err.Error())
		}
		m.Content = msg.Out
		m.viewport.SetContent(m.Content)
	}

	return m, tea.Batch(cmds...)
}

func (m StatusWindow) handleKeyPress(msg tea.KeyPressMsg) (StatusWindow, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "q":
		return m, tea.Quit
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m StatusWindow) View() string {
	color := lipgloss.White
	if m.Focus {
		color = lipgloss.Yellow
	}

	return m.Renderer.
		Width(m.Width).
		Height(m.Height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Render(m.viewport.View())
}

func (m StatusWindow) IsFocused() bool {
	return m.Focus
}
