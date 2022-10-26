package infrastructure

import (
	"context"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

type Input struct {
	AwsAccessKeyIDSecretID     dagger.SecretID
	AwsSecretAccessKeySecretID dagger.SecretID
	AwsSessionTokenSecretID    dagger.SecretID
	CatApiKeySecretID          dagger.SecretID
	BackendDistDirectoryID     dagger.DirectoryID
}

type BuildOutput struct{}

func Build(ctx context.Context, client *dagger.Client, src dagger.DirectoryID, ipt *Input) (*BuildOutput, error) {
	cont := setup(ctx, client, src, ipt)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"yarn", "plan"}})

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &BuildOutput{}, nil
}

type DeployOutput struct{}

func Deploy(ctx context.Context, client *dagger.Client, src dagger.DirectoryID, ipt *Input) (*DeployOutput, error) {
	cont := setup(ctx, client, src, ipt)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"yarn", "apply:auto-approve"}})

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &DeployOutput{}, nil
}

func newContainer(client *dagger.Client, src dagger.DirectoryID, ipt *Input) *dagger.Container {
	return util.
		NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/infrastructure").
		WithMountedDirectory("/app/backend/dist", ipt.BackendDistDirectoryID)
}

func setupSecrets(cont *dagger.Container, ipt *Input) *dagger.Container {
	return cont.
		WithSecretVariable("AWS_ACCESS_KEY_ID", ipt.AwsAccessKeyIDSecretID).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", ipt.AwsSecretAccessKeySecretID).
		WithSecretVariable("AWS_SESSION_TOKEN", ipt.AwsSessionTokenSecretID).
		WithSecretVariable("CAT_API_KEY", ipt.CatApiKeySecretID)
}

func setup(ctx context.Context, client *dagger.Client, src dagger.DirectoryID, ipt *Input) *dagger.Container {
	tfversion := "1.2.3"
	nodeversion := "16.x"

	cont := newContainer(client, src, ipt)
	cont = setupSecrets(cont, ipt)

	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "update", "-qq"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "install", "-y", "unzip"}})

	cont = util.SetupTerraform(cont, tfversion)
	cont = util.SetupNodeJS(cont, nodeversion)
	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"go", "get"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"yarn", "install", "--frozen-lockfile"}})

	return cont
}
