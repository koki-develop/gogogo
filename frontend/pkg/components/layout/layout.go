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

	root := util.WithClasses(elem.Div(h, m, f), "p-4", "relative")
	return root
}
