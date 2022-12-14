package main

import (
	"fmt"
	"os"

	"github.com/koki-develop/gogogo/cicd/pkg/backend"
	"github.com/koki-develop/gogogo/cicd/pkg/frontend"
	"github.com/koki-develop/gogogo/cicd/pkg/infrastructure"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("must pass workflow name")
		os.Exit(1)
	}

	workflow := os.Args[1]
	switch workflow {
	case "build", "deploy":
	default:
		fmt.Printf("unknown workflow: %s\n", workflow)
		os.Exit(1)
	}

	ctx, client, err := util.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	branch := os.Getenv("GITHUB_BRANCH")
	accessKeyID := util.Must(client.Host().EnvVariable("AWS_ACCESS_KEY_ID").Secret().ID(ctx))
	secretAccessKey := util.Must(client.Host().EnvVariable("AWS_SECRET_ACCESS_KEY").Secret().ID(ctx))
	sessionToken := util.Must(client.Host().EnvVariable("AWS_SESSION_TOKEN").Secret().ID(ctx))
	catApiKey := util.Must(client.Host().EnvVariable("CAT_API_KEY").Secret().ID(ctx))

	src := util.Must(util.Checkout(ctx, client, branch))

	bout := util.Must(backend.Build(ctx, client, src))

	_ = util.Must(infrastructure.Build(ctx, client, src, &infrastructure.BuildInput{
		AwsAccessKeyIDSecretID:     accessKeyID,
		AwsSecretAccessKeySecretID: secretAccessKey,
		AwsSessionTokenSecretID:    sessionToken,
		CatApiKeySecretID:          catApiKey,
		BackendDistDirectoryID:     bout.DistDirectoryID,
	}))

	fout := util.Must(frontend.Build(ctx, client, src))

	if workflow != "deploy" {
		return
	}

	_ = util.Must(infrastructure.Deploy(ctx, client, src, &infrastructure.DeployInput{
		AwsAccessKeyIDSecretID:     accessKeyID,
		AwsSecretAccessKeySecretID: secretAccessKey,
		AwsSessionTokenSecretID:    sessionToken,
		CatApiKeySecretID:          catApiKey,
		BackendDistDirectoryID:     bout.DistDirectoryID,
	}))

	_ = util.Must(frontend.Deploy(ctx, client, src, &frontend.DeployInput{
		AwsAccessKeyIDSecretID:     accessKeyID,
		AwsSecretAccessKeySecretID: secretAccessKey,
		AwsSessionTokenSecretID:    sessionToken,
		DistDirectoryID:            fout.DistDirectoryID,
	}))
}
