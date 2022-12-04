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
	case "all", "frontend", "backend", "infrastructure":
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
	if component == "all" || component == "backend" {
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
	if component == "all" || component == "infrastructure" {
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

		if workflow == "build" {
			// plan
			cont = cont.WithExec([]string{"yarn", "run", "cdktf", "plan"})
		}
		if workflow == "deploy" {
			// apply
			cont = cont.WithExec([]string{"yarn", "run", "cdktf", "apply", "--auto-approve"})
		}
		// exec pipeline
		if _, err := cont.ExitCode(ctx); err != nil {
			panic(err)
		}
	}

	// frontend
	if component == "all" || component == "frontend" {
		// checkout
		cont := client.Container().
			From("golang:1.19").
			WithMountedDirectory("/app", src).
			WithWorkdir("/app/frontend")

		if workflow == "build" {
			cont = cont.
				WithExec([]string{"go", "run", "./html"}).
				WithExec([]string{"bash", "-c", "cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./dist/wasm_exec.js"}).
				WithEnvVariable("GOOS", "js").
				WithEnvVariable("GOARCH", "wasm").
				WithExec([]string{"go", "build", "-o", "./dist/main.wasm"})
			fout := client.Directory().WithDirectory("frontend/dist", cont.Directory("./dist"))
			if _, err := fout.Export(ctx, root); err != nil {
				panic(err)
			}
		}
		if workflow == "deploy" {
			// install tools
			cont = util.SetupUnzip(cont)
			cont = util.SetupAWSCLI(cont)
			// setup environment variables and secrets
			cont = cont.
				WithEnvVariable("AWS_REGION", "us-east-1").
				WithSecretVariable("AWS_ACCESS_KEY_ID", accessKeyID).
				WithSecretVariable("AWS_SECRET_ACCESS_KEY", secretAccessKey).
				WithSecretVariable("AWS_SESSION_TOKEN", sessionToken)
			// deploy
			cont = cont.
				WithExec([]string{"aws", "s3", "cp", "./dist/index.html", "s3://gogogo-frontend-files/index.html"}).
				WithExec([]string{"aws", "s3", "cp", "./dist/wasm_exec.js", "s3://gogogo-frontend-files/wasm_exec.js"}).
				// 参考: https://stackoverflow.com/questions/51033550/how-to-manually-gzip-files-for-web-and-amazon-cloudfront
				WithExec([]string{"gzip", "./dist/main.wasm"}).
				WithExec([]string{"aws", "s3", "cp", "./dist/main.wasm.gz", "s3://gogogo-frontend-files/main.wasm", "--content-encoding", "gzip", "--content-type", "application/wasm"})
			// exec pipeline
			if _, err := cont.ExitCode(ctx); err != nil {
				panic(err)
			}
		}
	}
}
