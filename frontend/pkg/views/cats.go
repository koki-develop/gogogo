package views

import (
	"encoding/json"
	"net/http"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/frontend/pkg/components/layout"
	"github.com/koki-develop/gogogo/frontend/pkg/components/model"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

type CatsView struct {
	vecty.Core

	Cats      entities.Cats
	Loaded    bool
	Reloading bool
	Error     error
}

func (v *CatsView) Render() vecty.ComponentOrHTML {
	go func() {
		if v.Loaded && !v.Reloading {
			return
		}

		req, err := http.NewRequest(http.MethodGet, "https://dkasns0wq3.execute-api.us-east-1.amazonaws.com/prod/v1/cats", nil)
		if err != nil {
			v.Error = err
			v.Loaded = true
			v.Reloading = false
			vecty.Rerender(v)
			return
		}

		resp, err := new(http.Client).Do(req)
		if err != nil {
			v.Error = err
			v.Loaded = true
			v.Reloading = false
			vecty.Rerender(v)
			return
		}
		defer resp.Body.Close()

		var cats entities.Cats
		if err := json.NewDecoder(resp.Body).Decode(&cats); err != nil {
			v.Error = err
			v.Loaded = true
			v.Reloading = false
			vecty.Rerender(v)
			return
		}

		v.Cats = cats
		v.Loaded = true
		v.Reloading = false

		vecty.Rerender(v)
	}()

	if v.Error != nil {
		msg := elem.Paragraph(vecty.Text("Unexpected error occurred."))
		vecty.Markup(vecty.Class("text-center")).Apply(msg)
		return elem.Body(layout.New(msg))
	}

	// TODO: リファクタ
	var body vecty.MarkupOrChild
	if v.Loaded && !v.Reloading {
		body = elem.Div(
			elem.Div(
				vecty.Markup(
					vecty.Class("flex"),
					vecty.Class("justify-center"),
					vecty.Class("mb-2"),
				),
				elem.Button(
					vecty.Text("Reload"),
					vecty.Markup(
						event.Click(func(e *vecty.Event) {
							v.Reloading = true
							vecty.Rerender(v)
						}),
						vecty.Class("bg-blue-500"),
						vecty.Class("hover:bg-blue-700"),
						vecty.Class("text-white"),
						vecty.Class("font-bold"),
						vecty.Class("py-2"),
						vecty.Class("px-4"),
						vecty.Class("rounded"),
						vecty.Class("transition"),
					),
				),
			),
			model.NewCatImages(v.Cats),
		)
	} else {
		body = elem.Body(util.NewLoading())
	}

	return elem.Body(layout.New(body))
}

func NewCatsView() *CatsView {
	return &CatsView{}
}
