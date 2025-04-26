package windows

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/egit/pkg/ui/msgs"
)

type LogWindow struct {
	content string
}

type logCmdMsg msgs.GitCmdOutMsg

var runLogCmd = MakeGitRunner(func(out string, err error) tea.Msg {
	return logCmdMsg{Out: out, Err: err}
})

func (m LogWindow) Init() tea.Cmd {
	return runLogCmd("log", "--format='%h %s'")
}

func (m LogWindow) Update(msg tea.Msg) (WindowType, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		cmd = m.handleKeyPress(msg)
		cmds = append(cmds, cmd)

	case logCmdMsg:
		if msg.Err != nil {
			m.content = fmt.Sprintf("Error executing command: %v", msg.Err.Error())
		}

		m.content = msg.Out
	}

	return m, tea.Batch(cmds...)
}

func (m LogWindow) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.String() {
	case "q":
		return tea.Quit
	}

	return cmd
}

func (m LogWindow) Content() string {
	return m.content
}
