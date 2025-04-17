package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vieitesss/egit/pkg/cmd"
)

type Egit struct {
	renderer *lipgloss.Style
	width    int
	height   int
}

func NewEgit() tea.Model {
	return Egit{
		renderer: &lipgloss.Style{},
	}
}

func (m Egit) Init() tea.Cmd {
	return nil
}

func (m Egit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		cmd = m.handleKeyPress(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m Egit) View() string {
	var sections []string

	branchOut, _ := cmd.GitCmdOut("branch", "--show-current")
	statusOut, _ := cmd.GitCmdOut("status", "--short")
	logOut, _ := cmd.GitCmdOut("log", "--format='%h %s'")

	branch := m.renderer.
		Border(lipgloss.RoundedBorder()).
		Width(m.width).
		Height(m.height / 3).
		Render(nnl(branchOut))

	status := m.renderer.
		Border(lipgloss.RoundedBorder()).
		Width(m.width).
		Height(m.height / 3).
		Render(nnl(statusOut))

	log := m.renderer.
		Border(lipgloss.RoundedBorder()).
		SetString().
		Width(m.width).
		Height(m.height - lipgloss.Height(branch + status) - 1).
		Render(nnl(logOut))

	sections = append(sections, branch, status, log)

	return lipgloss.JoinVertical(0, sections...)
}

// Removes the trailing new lines of the string parameter
func nnl(s string) string {
	return strings.TrimRight(s, "\n")
}

func (m Egit) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	}

	return nil
}

func main() {
	pro := tea.NewProgram(NewEgit(), tea.WithAltScreen())

	if _, err := pro.Run(); err != nil {
		fmt.Errorf("Error initializing Egit: %v", err)
	}
}
