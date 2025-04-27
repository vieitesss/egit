package windows

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/egit/pkg/ui/msgs"
)

type StatusWindow struct {
	loaded bool
	content string
}

type statusCmdMsg msgs.GitCmdOutMsg

var runStatusCmd = MakeGitRunner(func(out string, err error) tea.Msg {
	return statusCmdMsg{Out: out, Err: err}
})

func (m StatusWindow) Init() tea.Cmd {
	return runStatusCmd("status", "--short")
}

func (m StatusWindow) Update(msg tea.Msg) (WindowType, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		cmd = m.handleKeyPress(msg)
		cmds = append(cmds, cmd)

	case statusCmdMsg:
		if !m.loaded {
			m.loaded = true
		}

		if msg.Err != nil {
			m.content = fmt.Sprintf("Error executing command: %v", msg.Err.Error())
		}

		m.content = msg.Out
	}

	return m, tea.Batch(cmds...)
}

func (m StatusWindow) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.String() {
	case "a":
		return runStatusCmd("status")
	case "r":
		return runStatusCmd("status", "--short")
	}

	return cmd
}

func (m StatusWindow) Content() string {
	return m.content
}

func (m StatusWindow) IsLoaded() bool {
	return m.loaded
}
