package windows

import (
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
	"github.com/vieitesss/egit/pkg/ui/msgs"
)

type WindowType interface {
	Init() tea.Cmd
	Update(tea.Msg) (WindowType, tea.Cmd)
	Content() string
}

type Window struct {
	comp.Component
	Type      WindowType
	viewport  viewport.Model
	MaxHeight int
	Loaded    bool
}

func NewWindow(ren *lipgloss.Style, winType WindowType, focus bool, maxHeight int) Window {
	return Window{
		Component: comp.Component{
			Renderer: ren,
			Focus:    focus,
		},
		Type:      winType,
		MaxHeight: maxHeight,
	}
}

func (m *Window) updateWindowSize(w, h int) {
	m.Width = w
	m.Height = h
	if m.MaxHeight > 0 && m.MaxHeight < h {
		m.Height = m.MaxHeight
	}

	m.viewport.SetWidth(m.Width - 2)
	m.viewport.SetHeight(m.Height - 2)
}

func (m Window) Init() tea.Cmd {
	m.viewport = viewport.New()

	return m.Type.Init()
}

func (m Window) Update(msg tea.Msg) (comp.ComponentI, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateWindowSize(msg.Width, msg.Height)
		return m, nil
	case msgs.GitCmdOutMsg:
		if !m.Focus && m.Loaded {
			return m, nil
		}
		m.Loaded = true
	}

	m.Type, cmd = m.Type.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport.SetContent(m.Type.Content())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Window) View() string {
	color := lipgloss.White
	if m.IsFocused() {
		color = lipgloss.Yellow
	}

	return m.Renderer.
		Width(m.Width).
		Height(m.Height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Render(m.viewport.View())
}

func (m Window) GetFixedHeight() int {
	return m.MaxHeight
}
