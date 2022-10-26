package frontend

import (
	"context"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

type BuildOutput struct {
	DistDirectoryID dagger.DirectoryID
}

func Build(ctx context.Context, client *dagger.Client, src dagger.DirectoryID) (*BuildOutput, error) {
	cont := newContainer(client, src)
	cont = util.SetupTask(cont)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"task", "build"}})

	d, err := cont.Directory("/app/frontend/dist").ID(ctx)
	if err != nil {
		panic(err)
	}

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &BuildOutput{DistDirectoryID: d}, nil
}

type DeployInput struct {
	AwsAccessKeyIDSecretID     dagger.SecretID
	AwsSecretAccessKeySecretID dagger.SecretID
	AwsSessionTokenSecretID    dagger.SecretID
	DistDirectoryID            dagger.DirectoryID
}

type DeployOutput struct{}

func Deploy(ctx context.Context, client *dagger.Client, src dagger.DirectoryID, ipt *DeployInput) (*DeployOutput, error) {
	cont := newContainer(client, src)
	cont = cont.WithMountedDirectory("/app/frontend/dist", ipt.DistDirectoryID)
	cont = setupSecrets(cont, &setupSecretsInput{
		AwsAccessKeyIDSecretID:     ipt.AwsAccessKeyIDSecretID,
		AwsSecretAccessKeySecretID: ipt.AwsSecretAccessKeySecretID,
		AwsSessionTokenSecretID:    ipt.AwsSessionTokenSecretID,
	})
	cont = util.SetupTask(cont)
	cont = util.SetupUnzip(cont)
	cont = util.SetupAWSCLI(cont)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"task", "deploy-only"}})

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &DeployOutput{}, nil
}

func newContainer(client *dagger.Client, src dagger.DirectoryID) *dagger.Container {
	return util.
		NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/frontend").
		WithEnvVariable("AWS_REGION", "us-east-1")
}

type setupSecretsInput struct {
	AwsAccessKeyIDSecretID     dagger.SecretID
	AwsSecretAccessKeySecretID dagger.SecretID
	AwsSessionTokenSecretID    dagger.SecretID
}

func setupSecrets(cont *dagger.Container, ipt *setupSecretsInput) *dagger.Container {
	return cont.
		WithSecretVariable("AWS_ACCESS_KEY_ID", ipt.AwsAccessKeyIDSecretID).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", ipt.AwsSecretAccessKeySecretID).
		WithSecretVariable("AWS_SESSION_TOKEN", ipt.AwsSessionTokenSecretID)
}
