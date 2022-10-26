package main

import (
	"log"

	"github.com/koki-develop/gogogo/cicd/pkg/backend"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

func main() {
	ctx, client, err := util.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	// backend
	_, err = backend.Build(ctx, client)
	if err != nil {
		log.Fatalln(err)
	}
}
