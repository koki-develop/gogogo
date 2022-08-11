package util

import "github.com/hexops/vecty"

func WithClasses(h *vecty.HTML, classes ...string) *vecty.HTML {
	as := []vecty.Applyer{}
	for _, c := range classes {
		as = append(as, vecty.Class(c))
	}

	vecty.Markup(as...).Apply(h)
	return h
}
