package windows

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/egit/pkg/ui/msgs"
)

type BranchWindow struct {
	loaded bool
	content string
}

type branchCmdMsg msgs.GitCmdOutMsg

var runBranchCmd = makeGitRunner(func(out string, err error) tea.Msg {
	return branchCmdMsg{Out: out, Err: err}
})

func (m BranchWindow) Init() tea.Cmd {
	return runBranchCmd("branch", "--show-current")
}

func (m BranchWindow) Update(msg tea.Msg) (WindowType, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		cmd = m.handleKeyPress(msg)
		cmds = append(cmds, cmd)

	case branchCmdMsg:
		if msg.Err != nil {
			m.content = fmt.Sprintf("Error executing command: %v", msg.Err.Error())
		}

		m.content = msg.Out
	}

	return m, tea.Batch(cmds...)
}

func (m BranchWindow) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.String() {
	}

	return cmd
}

func (m BranchWindow) Content() string {
	return m.content
}

func (m BranchWindow) IsLoaded() bool {
	return m.loaded
}
