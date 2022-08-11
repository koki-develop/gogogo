package layout

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

func newFooter() *vecty.HTML {
	copy := elem.Paragraph(
		vecty.Markup(
			vecty.UnsafeHTML("&copy;2022 Koki Sato"),
		),
	)

	f := util.WithClasses(elem.Footer(
		copy,
		util.NewAnchor("View on GitHub", "https://github.com/koki-develop/gogogo", true),
	), "text-center")

	return f
}
