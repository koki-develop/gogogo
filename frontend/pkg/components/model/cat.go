package model

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

type CatImage struct {
	vecty.Core

	Cat     *entities.Cat
	OnClick func(cat *entities.Cat)
}

func NewCatImage(cat *entities.Cat, onClick func(cat *entities.Cat)) *CatImage {
	return &CatImage{Cat: cat, OnClick: onClick}
}

func (c *CatImage) Render() vecty.ComponentOrHTML {
	img := util.WithClasses(util.NewImage(c.Cat.URL, "Cat"),
		"absolute", "w-full", "h-full", "top-0", "left-0", "object-cover",
	)

	card := elem.Div(vecty.Markup(vecty.Style("padding-top", "100%")), img)

	cardcontainer := util.WithClasses(elem.Div(vecty.Markup(event.Click(func(e *vecty.Event) { c.OnClick(c.Cat) })), card),
		"relative", "w-full", "rounded-full", "overflow-hidden", "shadow-lg", "border",
	)

	return cardcontainer
}

type CatImages struct {
	vecty.Core

	Cats       entities.Cats
	ShowingCat *entities.Cat
}

func NewCatImages(cats entities.Cats) *CatImages {
	return &CatImages{Cats: cats}
}

func (c *CatImages) onClickCat(cat *entities.Cat) {
	c.ShowingCat = cat
	vecty.Rerender(c)
}

func (c *CatImages) Render() vecty.ComponentOrHTML {
	var items []vecty.MarkupOrChild
	for _, cat := range c.Cats {
		img := NewCatImage(cat, c.onClickCat)
		item := util.WithClasses(elem.ListItem(img),
			"flex", "justify-center", "m-2",
		)

		items = append(items, item)
	}

	list := util.WithClasses(elem.UnorderedList(items...),
		"grid", "grid-cols-2", "sm:grid-cols-3", "md:grid-cols-4",
	)

	body := []vecty.MarkupOrChild{list}
	if c.ShowingCat != nil {
		body = append(body, elem.Div(
			vecty.Markup(
				event.Click(func(e *vecty.Event) {
					c.ShowingCat = nil
					vecty.Rerender(c)
				}),
				vecty.Style("background-color", "rgba(0, 0, 0, 0.5)"),
				vecty.Class("fixed"),
				vecty.Class("flex"),
				vecty.Class("items-center"),
				vecty.Class("justify-center"),
				vecty.Class("w-full"),
				vecty.Class("h-full"),
				vecty.Class("left-0"),
				vecty.Class("top-0"),
				vecty.Class("z-50"),
			),

			elem.Image(
				vecty.Markup(
					vecty.Attribute("src", c.ShowingCat.URL),
					vecty.Class("object-contain"),
					vecty.Style("min-width", "60%"),
					vecty.Style("min-height", "60%"),
					vecty.Style("max-width", "80%"),
					vecty.Style("max-height", "80%"),
				),
			),
		))
	}

	return elem.Div(body...)
}
