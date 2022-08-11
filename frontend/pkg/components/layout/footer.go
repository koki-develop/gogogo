package layout

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

func newFooter() *vecty.HTML {
	copy := elem.Paragraph()
	vecty.Markup(vecty.UnsafeHTML("&copy;2022 Koki Sato")).Apply(copy)

	f := elem.Footer(
		copy,
		util.NewAnchor("View on GitHub", "https://github.com/koki-develop/gogogo", true),
	)
	vecty.Markup(vecty.Class("text-center")).Apply(f)

	return f
}
