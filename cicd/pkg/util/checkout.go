package util

import (
	"context"

	"dagger.io/dagger"
)

func Checkout(ctx context.Context, client *dagger.Client, img string) (*dagger.Container, error) {
	repoUrl := "https://github.com/koki-develop/gogogo.git"
	branch := "main"

	repo := client.Git(repoUrl)
	src, err := repo.Branch(branch).Tree().ID(ctx)
	if err != nil {
		return nil, err
	}

	cont := client.Container().From(img)
	cont = cont.WithMountedDirectory("/app", src)

	return cont, nil
}
