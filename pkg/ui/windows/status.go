package windows

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/egit/pkg/ui/msgs"
	"github.com/vieitesss/egit/pkg/cmd"
)

type StatusWindow struct {
	content string
}

func (m StatusWindow) Init() tea.Cmd {
	return runGitCmdOut("status")
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

	case msgs.GitCmdOutMsg:
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
	case "q":
		return tea.Quit
	}

	return cmd
}

func (m StatusWindow) Content() string {
	return m.content
}

func runGitCmdOut(c ...string) tea.Cmd {
	return func() tea.Msg {
		out, err := cmd.GitCmdOut(c...)
		return msgs.GitCmdOutMsg{
			Out: out,
			Err: err,
		}
	}
}
