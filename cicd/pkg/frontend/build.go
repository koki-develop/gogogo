package frontend

import (
	"context"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

type BuildOutput struct{}

func Build(ctx context.Context, client *dagger.Client, src dagger.DirectoryID) (*BuildOutput, error) {
	// initialize
	cont := util.
		NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/frontend")

	// setup go-task
	cont = util.SetupTask(cont)

	// build
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"task", "build"}})

	// run
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
	// initialize
	cont := util.
		NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/frontend")

	// install unzip
	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "update", "-qq"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "install", "-y", "unzip"}})

	// setup go-task
	cont = util.SetupTask(cont)

	// setup awscli
	cont = util.SetupAWSCLI(cont)

	// build
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"task", "build"}})

	// run
	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}

	return &DeployOutput{}, nil
}
