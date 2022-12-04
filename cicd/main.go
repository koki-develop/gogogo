package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("must pass component and workflow")
		os.Exit(1)
	}

	component := os.Args[1]
	switch component {
	case "frontend", "backend", "infrastructure":
	default:
		fmt.Printf("unknown component: %s\n", component)
		os.Exit(1)
	}
	workflow := os.Args[2]
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
	accessKeyID := client.Host().EnvVariable("AWS_ACCESS_KEY_ID").Secret()
	secretAccessKey := client.Host().EnvVariable("AWS_SECRET_ACCESS_KEY").Secret()
	sessionToken := client.Host().EnvVariable("AWS_SESSION_TOKEN").Secret()
	catApiKey := client.Host().EnvVariable("CAT_API_KEY").Secret()

	// get src
	root, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}
	src := client.Host().Directory(root)

	// backend
	if component == "backend" {
		if workflow == "build" {
			// checkout
			cont := client.Container().
				From("golang:1.19").
				WithMountedDirectory("/app", src).
				WithWorkdir("/app/backend")
			// setup environment variables
			cont = cont.
				WithEnvVariable("GOARCH", "amd64").
				WithEnvVariable("GOOS", "linux")
			// build
			cont = cont.
				WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", "dist/api", "./pkg/handlers/api/lambda"}).
				WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", "dist/updatecats", "./pkg/handlers/updatecats"})
			// exec pipeline
			bout := client.Directory().WithDirectory("backend/dist", cont.Directory("./dist"))
			if _, err := bout.Export(ctx, root); err != nil {
				panic(err)
			}
		}
	}

	// infrastructure
	if component == "infrastructure" {
		if workflow == "build" {
			// checkout
			cont := client.Container().
				From("golang:1.19").
				WithMountedDirectory("/app", src).
				WithWorkdir("/app/infrastructure")
			// setup environment variables and secrets
			cont = cont.
				WithEnvVariable("CI", "true").
				WithSecretVariable("AWS_ACCESS_KEY_ID", accessKeyID).
				WithSecretVariable("AWS_SECRET_ACCESS_KEY", secretAccessKey).
				WithSecretVariable("AWS_SESSION_TOKEN", sessionToken).
				WithSecretVariable("CAT_API_KEY", catApiKey)
			// install tools
			cont = util.SetupUnzip(cont)
			cont = util.SetupTerraform(cont, "1.2.3")
			cont = util.SetupNodeJS(cont, "16.x")
			// install dependencies
			cont = cont.
				WithExec([]string{"go", "get"}).
				WithExec([]string{"yarn", "install", "--frozen-lockfile"})
			// plan
			cont = cont.WithExec([]string{"yarn", "run", "cdktf", "plan"})
			// exec pipeline
			if _, err := cont.ExitCode(ctx); err != nil {
				panic(err)
			}
		}

	}

	// frontend
	if component == "frontend" {
		if workflow == "build" {
			// checkout
			cont := client.Container().
				From("golang:1.19").
				WithMountedDirectory("/app", src).
				WithWorkdir("/app/frontend")
			// setup environment variables
			cont = cont.
				WithEnvVariable("AWS_REGION", "us-east-1")
			// build
			cont = cont.
				WithExec([]string{"go", "run", "./html"}).
				WithExec([]string{"bash", "-c", "cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./dist/wasm_exec.js"}).
				WithEnvVariable("GOOS", "js").
				WithEnvVariable("GOARCH", "wasm").
				WithExec([]string{"go", "build", "-o", "./dist/main.wasm"})
			// exec pipeline
			if _, err := cont.ExitCode(ctx); err != nil {
				panic(err)
			}
		}
	}
}
