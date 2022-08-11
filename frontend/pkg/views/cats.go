package views

import (
	"encoding/json"
	"net/http"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/frontend/pkg/components/layout"
	"github.com/koki-develop/gogogo/frontend/pkg/components/model"
)

type CatsView struct {
	vecty.Core
}

func (v *CatsView) Render() vecty.ComponentOrHTML {
	req, _ := http.NewRequest(http.MethodGet, "https://dkasns0wq3.execute-api.us-east-1.amazonaws.com/prod/v1/cats", nil)
	resp, _ := new(http.Client).Do(req)
	defer resp.Body.Close()

	var cats entities.Cats
	json.NewDecoder(resp.Body).Decode(&cats)

	imgs := model.NewCatImages(cats)

	return elem.Body(layout.New(
		imgs,
	))
}

func NewCatsView() *CatsView {
	return &CatsView{}
}
