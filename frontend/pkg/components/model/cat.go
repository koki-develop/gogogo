package model

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

type CatImage struct {
	vecty.Core

	Cat *entities.Cat
}

func NewCatImage(cat *entities.Cat) *CatImage {
	return &CatImage{Cat: cat}
}

func (c *CatImage) Render() vecty.ComponentOrHTML {
	img := util.WithClasses(util.NewImage(c.Cat.URL, "Cat"),
		"absolute", "w-full", "h-full", "top-0", "left-0", "object-cover",
	)

	card := elem.Div(img)
	vecty.Markup(
		vecty.Style("padding-top", "100%"),
	).Apply(card)

	cardcontainer := util.WithClasses(elem.Div(card),
		"relative", "w-full", "rounded-full", "overflow-hidden", "shadow-lg", "border",
	)

	return cardcontainer
}

type CatImages struct {
	vecty.Core

	Cats entities.Cats
}

func NewCatImages(cats entities.Cats) *CatImages {
	return &CatImages{Cats: cats}
}

func (c *CatImages) Render() vecty.ComponentOrHTML {
	var items []vecty.MarkupOrChild
	for _, cat := range c.Cats {
		img := NewCatImage(cat)
		item := util.WithClasses(elem.ListItem(img),
			"flex", "justify-center", "m-2",
		)

		items = append(items, item)
	}

	list := util.WithClasses(elem.UnorderedList(items...),
		"grid", "grid-cols-2", "sm:grid-cols-3", "md:grid-cols-4",
	)

	return list
}
