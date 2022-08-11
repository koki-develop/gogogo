package util

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func NewAnchor(text, href string, external bool) *vecty.HTML {
	a := elem.Anchor(vecty.Text(text))

	attrs := []vecty.Applyer{vecty.Attribute("href", href)}
	if external {
		attrs = append(attrs, vecty.Attribute("target", "_blank"))
		attrs = append(attrs, vecty.Attribute("rel", "noopener noreferrer"))
	}

	vecty.Markup(attrs...).Apply(a)

	return a
}
