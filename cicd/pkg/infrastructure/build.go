package infrastructure

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

type BuildInput struct {
	AwsAccessKeyIDSecretID     dagger.SecretID
	AwsSecretAccessKeySecretID dagger.SecretID
	AwsSessionTokenSecretID    dagger.SecretID
	CatApiKeySecretID          dagger.SecretID
	BackendDistDirectoryID     dagger.DirectoryID
}

type BuildOutput struct{}

func Build(ctx context.Context, client *dagger.Client, src dagger.DirectoryID, ipt *BuildInput) (*BuildOutput, error) {
	tfversion := "1.2.3"
	nodeversion := "16.x"

	// initialize
	cont := util.NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/infrastructure").
		WithMountedDirectory("/app/backend/dist", ipt.BackendDistDirectoryID)

	// secrets
	cont = cont.WithSecretVariable("AWS_ACCESS_KEY_ID", ipt.AwsAccessKeyIDSecretID).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", ipt.AwsSecretAccessKeySecretID).
		WithSecretVariable("AWS_SESSION_TOKEN", ipt.AwsSessionTokenSecretID).
		WithSecretVariable("CAT_API_KEY", ipt.CatApiKeySecretID)

	// install unzip
	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "update", "-qq"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "install", "-y", "unzip"}})

	// install terraform
	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{
			"wget",
			fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_linux_amd64.zip", tfversion, tfversion),
			"-O",
			"/tmp/terraform.zip",
		}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"unzip", "/tmp/terraform.zip", "-d", "/usr/bin"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"chmod", "+x", "/usr/bin/terraform"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"terraform", "--version"}})

	// install nodejs
	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"bash", "-c", fmt.Sprintf("curl -fsSL https://deb.nodesource.com/setup_%s | bash -", nodeversion)}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "install", "-y", "nodejs"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"npm", "install", "-g", "yarn"}})

	// install dependencies
	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"go", "get"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"yarn", "install", "--frozen-lockfile"}})

	// plan
	cont = cont.Exec(dagger.ContainerExecOpts{
		Args: []string{"yarn", "plan"},
	})

	// run
	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}

	return &BuildOutput{}, nil
}
