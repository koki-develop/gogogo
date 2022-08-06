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
		vecty.Attribute("width", 120),
	).Apply(img)
	return img
}

func NewCatImages(cats entities.Cats) vecty.MarkupOrChild {
	var items []vecty.MarkupOrChild
	for _, cat := range cats {
		img := NewCatImage(cat)
		items = append(items, elem.ListItem(img))
	}

	return elem.UnorderedList(items...)
}
