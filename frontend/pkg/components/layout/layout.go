package layout

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

func New(children ...vecty.MarkupOrChild) vecty.MarkupOrChild {
	h := newHeader()

	m := util.WithClasses(elem.Main(children...), "container", "my-4", "mx-auto")

	f := newFooter()

	root := elem.Div(vecty.Markup(
		vecty.Class("p-4"),
		vecty.Class("relative"),
	), h, m, f)
	return root
}
