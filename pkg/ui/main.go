package ui

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
	cont "github.com/vieitesss/egit/pkg/ui/containers"
	win "github.com/vieitesss/egit/pkg/ui/windows"
)

type Egit struct {
	renderer *lipgloss.Style
	width    int
	height   int
	layout   comp.Component
}

func DefaultLayout(ren *lipgloss.Style) comp.Component {
	w := win.Window{Renderer: ren}
	f := win.Window{Renderer: ren, Focus: true}
	return cont.NewColumn(
		ren,
		cont.NewRow(
			ren,
			cont.NewColumn(
				ren,
				win.StatusWindow{f},
				win.StatusWindow{w},
				win.StatusWindow{w},
			),
			cont.NewColumn(
				ren,
				win.StatusWindow{w},
				win.StatusWindow{w},
			),
		),
		win.StatusWindow{w},
	)
}

func NewEgit() tea.Model {
	ren := &lipgloss.Style{}
	return Egit{
		renderer: ren,
		layout:   DefaultLayout(ren),
	}
}

func (m Egit) Init() tea.Cmd {
	return m.layout.Init()
}

func (m Egit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.layout, cmd = m.layout.Update(msg)
	return m, cmd
}

func (m Egit) View() string {
	return m.layout.View()
}
