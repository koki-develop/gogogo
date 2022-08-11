package main

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/koki-develop/gogogo/backend/pkg/infrastructures/catapi"
	"github.com/koki-develop/gogogo/backend/pkg/infrastructures/s3"
)

func main() {
	lambda.Start(handler)
}

func handler() error {
	apiKey := os.Getenv("CAT_API_KEY")
	catcl := catapi.New(apiKey)

	s3cl := s3.New()

	cats, err := catcl.Search()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cats)
	if err != nil {
		return err
	}

	if err := s3cl.Upload("gogogo-cats", "cats.json", "application/json", bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
