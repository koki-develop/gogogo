package layout

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

func newHeader() *vecty.HTML {
	ttl := util.WithClasses(elem.Heading1(vecty.Text("GoGoGo")), "text-3xl")

	catapilink := util.WithClasses(util.NewAnchor("The Cat API", "https://thecatapi.com/", true), "text-blue-500")
	subttl := elem.Paragraph(vecty.Text("powered by "), catapilink)

	h := util.WithClasses(elem.Header(ttl, subttl), "text-center")

	return h
}
