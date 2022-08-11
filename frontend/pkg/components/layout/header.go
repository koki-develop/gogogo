package layout

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

func newHeader() *vecty.HTML {
	ttl := elem.Heading1(vecty.Text("GoGoGo"))
	vecty.Markup(vecty.Class("text-3xl")).Apply(ttl)

	catapilink := util.NewAnchor("The Cat API", "https://thecatapi.com/", true)
	vecty.Markup(vecty.Class("text-blue-500")).Apply(catapilink)
	subttl := elem.Paragraph(vecty.Text("powered by "), catapilink)

	h := elem.Header(ttl, subttl)
	vecty.Markup(vecty.Class("text-center")).Apply(h)

	return h
}
