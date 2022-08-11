package catapi

import (
	"encoding/json"
	"net/http"

	"github.com/koki-develop/gogogo/backend/pkg/entities"
)

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (cl *Client) Search() (entities.Cats, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.thecatapi.com/v1/images/search?limit=100", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Api-Key", cl.apiKey)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cats entities.Cats
	if err := json.NewDecoder(resp.Body).Decode(&cats); err != nil {
		return nil, err
	}

	return cats, nil
}
