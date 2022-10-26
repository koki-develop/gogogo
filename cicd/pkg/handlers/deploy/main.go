package main

import (
	"log"
	"os"

	"github.com/koki-develop/gogogo/cicd/pkg/backend"
	"github.com/koki-develop/gogogo/cicd/pkg/infrastructure"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

func main() {
	ctx, client, err := util.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	branch := os.Getenv("GITHUB_BRANCH")
	accessKeyID := util.Must(client.Host().EnvVariable("AWS_ACCESS_KEY_ID").Secret().ID(ctx))
	secretAccessKey := util.Must(client.Host().EnvVariable("AWS_SECRET_ACCESS_KEY").Secret().ID(ctx))
	sessionToken := util.Must(client.Host().EnvVariable("AWS_SESSION_TOKEN").Secret().ID(ctx))
	catApiKey := util.Must(client.Host().EnvVariable("CAT_API_KEY").Secret().ID(ctx))

	src, err := util.Checkout(ctx, client, branch)
	if err != nil {
		log.Fatalln(err)
	}

	// backend
	bout, err := backend.Build(ctx, client, src)
	if err != nil {
		log.Fatalln(err)
	}

	// infrastructure
	_, err = infrastructure.Deploy(ctx, client, src, &infrastructure.Input{
		AwsAccessKeyIDSecretID:     accessKeyID,
		AwsSecretAccessKeySecretID: secretAccessKey,
		AwsSessionTokenSecretID:    sessionToken,
		CatApiKeySecretID:          catApiKey,
		BackendDistDirectoryID:     bout.DistDirectoryID,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
