package model

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

func NewCatImage(cat *entities.Cat) vecty.MarkupOrChild {
	img := util.NewImage(cat.URL, "Cat")
	vecty.Markup(
		vecty.Class("absolute"),
		vecty.Class("w-full"),
		vecty.Class("h-full"),
		vecty.Class("top-0"),
		vecty.Class("left-0"),
		vecty.Class("object-cover"),
	).Apply(img)

	card := elem.Div(img)
	vecty.Markup(
		vecty.Style("padding-top", "100%"),
	).Apply(card)

	cardcontainer := elem.Div(card)
	vecty.Markup(
		vecty.Class("relative"),
		vecty.Class("w-full"),
		vecty.Class("rounded-full"),
		vecty.Class("overflow-hidden"),
		vecty.Class("shadow-lg"),
		vecty.Class("border"),
	).Apply(cardcontainer)

	return cardcontainer
}

func NewCatImages(cats entities.Cats) vecty.MarkupOrChild {
	var items []vecty.MarkupOrChild
	for _, cat := range cats {
		img := NewCatImage(cat)
		item := elem.ListItem(img)
		vecty.Markup(
			vecty.Class("flex"),
			vecty.Class("justify-center"),
			vecty.Class("m-2"),
		).Apply(item)

		items = append(items, item)
	}

	list := elem.UnorderedList(items...)
	vecty.Markup(
		vecty.Class("grid"),
		vecty.Class("grid-cols-2"),
		vecty.Class("sm:grid-cols-3"),
		vecty.Class("md:grid-cols-4"),
	).Apply(list)

	return list
}
