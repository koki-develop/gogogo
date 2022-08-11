package views

import (
	"encoding/json"
	"net/http"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/frontend/pkg/components/layout"
	"github.com/koki-develop/gogogo/frontend/pkg/components/model"
	"github.com/koki-develop/gogogo/frontend/pkg/components/util"
)

type CatsView struct {
	vecty.Core

	Cats   entities.Cats
	Loaded bool
	Error  error
}

func (v *CatsView) Render() vecty.ComponentOrHTML {
	go func() {
		if v.Loaded {
			return
		}

		req, err := http.NewRequest(http.MethodGet, "https://dkasns0wq3.execute-api.us-east-1.amazonaws.com/prod/v1/cats", nil)
		if err != nil {
			v.Error = err
			v.Loaded = true
			vecty.Rerender(v)
			return
		}

		resp, err := new(http.Client).Do(req)
		if err != nil {
			v.Error = err
			v.Loaded = true
			vecty.Rerender(v)
			return
		}
		defer resp.Body.Close()

		var cats entities.Cats
		if err := json.NewDecoder(resp.Body).Decode(&cats); err != nil {
			v.Error = err
			v.Loaded = true
			vecty.Rerender(v)
			return
		}

		v.Cats = cats
		v.Loaded = true

		vecty.Rerender(v)
	}()

	if v.Error != nil {
		msg := elem.Paragraph(vecty.Text("Unexpected error occurred."))
		vecty.Markup(vecty.Class("text-center")).Apply(msg)
		return elem.Body(layout.New(msg))
	}

	var body vecty.MarkupOrChild
	if v.Loaded {
		body = model.NewCatImages(v.Cats)
	} else {
		body = elem.Body(util.NewLoading())
	}

	return elem.Body(layout.New(body))
}

func NewCatsView() *CatsView {
	return &CatsView{}
}
