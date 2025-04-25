package containers

import (
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
)

type Row struct{}

func NewRow(ren *lipgloss.Style, cmps ...comp.ComponentI) Container {
	cont := newContainer(ren, cmps...)
	cont.Type = Row{}
	return cont
}

func (c Row) Joiner(vs []string) string {
	return lipgloss.JoinHorizontal(0, vs...)
}

func (c Row) Dimension(count, w, h int) (int, int) {
	return w / count, h
}
