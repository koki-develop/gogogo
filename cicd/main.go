package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
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

	// initialize client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// get secrets
	// accessKeyID := client.Host().EnvVariable("AWS_ACCESS_KEY_ID").Secret()
	// secretAccessKey := client.Host().EnvVariable("AWS_SECRET_ACCESS_KEY").Secret()
	// sessionToken := client.Host().EnvVariable("AWS_SESSION_TOKEN").Secret()
	// catApiKey := client.Host().EnvVariable("CAT_API_KEY").Secret()

	// get src
	root, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}
	src := client.Host().Directory(root)

	// backend
	{
		// checkout
		cont := client.Container().From("golang:1.19").
			WithMountedDirectory("/app", src).
			WithWorkdir("/app/backend")
		// setup environment variablees
		cont = cont.WithEnvVariable("GOARCH", "amd64").
			WithEnvVariable("GOOS", "linux")
		// build
		cont = cont.
			WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", "dist/api", "./pkg/handlers/api/lambda"}).
			WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", "dist/updatecats", "./pkg/handlers/updatecats"})
		// run pipeline
		if _, err := cont.ExitCode(ctx); err != nil {
			panic(err)
		}
	}
}
