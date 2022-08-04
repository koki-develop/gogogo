package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type PageView struct {
	vecty.Core
}

func (v *PageView) Render() vecty.ComponentOrHTML {
	vecty.SetTitle("GoGoGo")

	return elem.Body(
		elem.Heading1(vecty.Text("GoGoGo")),
		elem.Paragraph(vecty.Text("Hello World")),
	)
}

func main() {
	v := &PageView{}

	vecty.RenderBody(v)
}
