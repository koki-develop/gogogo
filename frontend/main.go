package main

import (
	"github.com/hexops/vecty"
	"github.com/koki-develop/gogogo/frontend/pkg/views"
)

func main() {
	v := views.NewCatsView()
	vecty.RenderBody(v)
}
