package util

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func NewImage(src, alt string) *vecty.HTML {
	img := elem.Image(
		vecty.Markup(
			vecty.Attribute("src", src),
			vecty.Attribute("alt", alt),
		),
	)

	return img
}
