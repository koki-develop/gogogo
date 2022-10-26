package util

import (
	"context"
	"os"

	"dagger.io/dagger"
)

func NewClient() (context.Context, *dagger.Client, error) {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return nil, nil, err
	}

	return ctx, client, nil
}
