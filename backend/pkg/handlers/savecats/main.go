package main

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/koki-develop/gogogo/backend/pkg/infrastructures/catsapi"
	"github.com/koki-develop/gogogo/backend/pkg/infrastructures/s3"
)

func main() {
	apiKey := os.Getenv("CATS_API_KEY")

	catscl := catsapi.New(apiKey)
	s3cl := s3.New()

	cats, err := catscl.Search()
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(cats)
	if err != nil {
		panic(err)
	}

	if err := s3cl.Upload("gogogo-cats", "cats.json", "application/json", bytes.NewBuffer(data)); err != nil {
		panic(err)
	}
}
