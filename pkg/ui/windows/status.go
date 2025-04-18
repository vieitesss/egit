package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vieitesss/egit/pkg/cmd"
	"github.com/vieitesss/egit/pkg/ui/msgs"
	comp "github.com/vieitesss/egit/pkg/ui/component"
)

type Window comp.Component

type StatusWindow struct {
	Window
	Renderer *lipgloss.Style
	Width    int
	Height   int
	Focus    bool
	Content  string
}

func (m StatusWindow) Init() tea.Cmd {
	return func() tea.Msg {
		out, err := cmd.GitCmdOut("status", "--porcelain")
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

	case tea.KeyPressMsg:
		m, cmd = m.handleKeyPress(msg.String())
		cmds = append(cmds, cmd)

	case msgs.StatusCmdMsg:
		if msg.Err != nil {
			m.Content = fmt.Sprintf("Error executing command: %v", msg.Err.Error())
		}
		m.Content = msg.Out
	}

	return m, tea.Batch(cmds...)
}

func (m StatusWindow) handleKeyPress(s string) (StatusWindow, tea.Cmd) {
	switch s {
	case "q":
		return m, tea.Quit
	}

	return m, nil
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
		Render(m.Content)
}

func (m StatusWindow) IsFocused() bool {
	return m.Focus
}
