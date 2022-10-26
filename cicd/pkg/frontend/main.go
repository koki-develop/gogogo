package frontend

import (
	"context"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

type BuildOutput struct{}

func Build(ctx context.Context, client *dagger.Client, src dagger.DirectoryID) (*BuildOutput, error) {
	cont := newContainer(client, src)
	cont = util.SetupTask(cont)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"task", "build"}})

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &BuildOutput{}, nil
}

type DeployInput struct {
	AwsAccessKeyIDSecretID     dagger.SecretID
	AwsSecretAccessKeySecretID dagger.SecretID
	AwsSessionTokenSecretID    dagger.SecretID
}

type DeployOutput struct{}

func Deploy(ctx context.Context, client *dagger.Client, src dagger.DirectoryID, ipt *DeployInput) (*DeployOutput, error) {
	cont := newContainer(client, src)
	cont = setupSecrets(cont, ipt)

	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "update", "-qq"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "install", "-y", "unzip"}})

	cont = util.SetupTask(cont)
	cont = util.SetupAWSCLI(cont)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"task", "deploy"}})

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &DeployOutput{}, nil
}

func newContainer(client *dagger.Client, src dagger.DirectoryID) *dagger.Container {
	return util.
		NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/frontend")
}

func setupSecrets(cont *dagger.Container, ipt *DeployInput) *dagger.Container {
	return cont.
		WithSecretVariable("AWS_ACCESS_KEY_ID", ipt.AwsAccessKeyIDSecretID).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", ipt.AwsSecretAccessKeySecretID).
		WithSecretVariable("AWS_SESSION_TOKEN", ipt.AwsSessionTokenSecretID)
}
