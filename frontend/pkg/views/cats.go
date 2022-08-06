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

	ttl := elem.Heading1(
		vecty.Text("GoGoGo"),
	)
	vecty.Markup(
		vecty.Class("text-center"),
		vecty.Class("text-3xl"),
	).Apply(ttl)

	catsapilink := elem.Anchor(vecty.Text("The Cat API"))
	vecty.Markup(
		vecty.Class("text-blue-500"),
		vecty.Attribute("href", "https://thecatapi.com/"),
		vecty.Attribute("target", "_blank"),
		vecty.Attribute("rel", "noopener noreferrer"),
	).Apply(catsapilink)

	poweredby := elem.Paragraph(vecty.Text("powered by "), catsapilink)
	vecty.Markup(
		vecty.Class("text-center"),
		vecty.Class("mb-4"),
	).Apply(poweredby)

	imgs := components.NewCatImages(cats)

	repolink := elem.Anchor(vecty.Text("View on GitHub"))
	vecty.Markup(
		vecty.Attribute("href", "https://github.com/koki-develop/gogogo"),
		vecty.Attribute("target", "_blank"),
		vecty.Attribute("rel", "noopener noreferrer"),
	).Apply(repolink)

	footer := elem.Div(repolink)
	vecty.Markup(vecty.Class("text-center"), vecty.Class("my-12")).Apply(footer)

	container := elem.Div(
		ttl,
		poweredby,
		imgs,
		footer,
	)
	vecty.Markup(
		vecty.Class("container"),
	).Apply(container)

	root := elem.Div(container)
	vecty.Markup(
		vecty.Class("flex"),
		vecty.Class("justify-center"),
		vecty.Class("p-4"),
		vecty.Class("sm:px-16"),
	).Apply(root)

	return elem.Body(root)
}

func NewCatsView() *CatsView {
	return &CatsView{}
}
