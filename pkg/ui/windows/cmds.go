package windows

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vieitesss/egit/pkg/cmd"
)

func MakeGitRunner(wrap func(out string, err error) tea.Msg) func(args ...string) tea.Cmd {
	return func(args ...string) tea.Cmd {
		out, err := cmd.GitCmdOut(args...)
		return func() tea.Msg {
			return wrap(out, err)
		}
	}
}
