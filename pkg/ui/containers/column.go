package containers

import (
	"github.com/charmbracelet/lipgloss/v2"
	comp "github.com/vieitesss/egit/pkg/ui/component"
)

type Column struct{}

func NewColumn(ren *lipgloss.Style, cmps ...comp.ComponentI) Container {
	cont := newContainer(ren, cmps...)
	cont.Type = Column{}
	return cont
}

func (c Column) Joiner(vs []string) string {
	return lipgloss.JoinVertical(0, vs...)
}

func (c Column) Dimension(count, w, h int) (int, int) {
	return w, h / count
}
