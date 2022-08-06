package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
)

func NewCatImage(cat *entities.Cat) vecty.MarkupOrChild {
	img := elem.Image()
	vecty.Markup(
		vecty.Attribute("src", cat.URL),
		vecty.Class("w-full"),
		vecty.Class("object-cover"),
	).Apply(img)
	return img
}

func NewCatImages(cats entities.Cats) vecty.MarkupOrChild {
	var items []vecty.MarkupOrChild
	for _, cat := range cats {
		img := NewCatImage(cat)
		item := elem.ListItem(img)
		vecty.Markup(
			vecty.Class("flex"),
			vecty.Class("justify-center"),
		).Apply(item)

		items = append(items, item)
	}

	list := elem.UnorderedList(items...)
	vecty.Markup(
		vecty.Class("grid"),
		vecty.Class("grid-cols-2"),
		vecty.Class("md:grid-cols-4"),
	).Apply(list)

	return list
}
