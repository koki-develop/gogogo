package main

import (
	"log"
	"os"

	"github.com/koki-develop/gogogo/cicd/pkg/backend"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

func main() {
	branch := os.Getenv("GITHUB_BRANCH")

	ctx, client, err := util.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	src, err := util.Checkout(ctx, client, branch)
	if err != nil {
		log.Fatalln(err)
	}

	// backend
	_, err = backend.Build(ctx, client, src)
	if err != nil {
		log.Fatalln(err)
	}
}
