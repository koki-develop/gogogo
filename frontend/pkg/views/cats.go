package views

import (
	"encoding/json"
	"net/http"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/frontend/pkg/components"
)

type CatsView struct {
	vecty.Core
}

func (v *CatsView) Render() vecty.ComponentOrHTML {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/cats", nil)
	resp, _ := new(http.Client).Do(req)
	defer resp.Body.Close()

	var cats entities.Cats
	json.NewDecoder(resp.Body).Decode(&cats)

	imgs := components.NewCatImages(cats)

	container := elem.Div(imgs)
	vecty.Markup(
		vecty.Class("container"),
	).Apply(container)

	root := elem.Div(container)
	vecty.Markup(
		vecty.Class("flex"),
		vecty.Class("justify-center"),
		vecty.Class("px-4"),
		vecty.Class("sm:px-16"),
	).Apply(root)

	return elem.Body(root)
}

func NewCatsView() *CatsView {
	return &CatsView{}
}
