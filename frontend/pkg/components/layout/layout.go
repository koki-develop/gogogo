package layout

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func New(children ...vecty.MarkupOrChild) vecty.MarkupOrChild {
	h := newHeader()

	m := elem.Main(children...)
	vecty.Markup(vecty.Class("my-4"), vecty.Class("container"), vecty.Class("mx-auto")).Apply(m)

	f := newFooter()

	root := elem.Div(h, m, f)
	vecty.Markup(vecty.Class("p-4")).Apply(root)
	return root
}
