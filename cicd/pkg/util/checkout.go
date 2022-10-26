package util

import (
	"context"

	"dagger.io/dagger"
)

func Checkout(ctx context.Context, client *dagger.Client, branch string) (dagger.DirectoryID, error) {
	repoUrl := "https://github.com/koki-develop/gogogo.git"

	repo := client.Git(repoUrl)
	src, err := repo.Branch(branch).Tree().ID(ctx)
	if err != nil {
		return "", err
	}

	return src, nil
}
