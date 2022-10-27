package infrastructure

import (
	"context"

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
	cont := newContainer(client, src, &newContainerInput{BackendDistDirectoryID: ipt.BackendDistDirectoryID})
	cont = setupSecrets(cont, &setupSecretsInput{
		AwsAccessKeyIDSecretID:     ipt.AwsAccessKeyIDSecretID,
		AwsSecretAccessKeySecretID: ipt.AwsSecretAccessKeySecretID,
		AwsSessionTokenSecretID:    ipt.AwsSessionTokenSecretID,
		CatApiKeySecretID:          ipt.CatApiKeySecretID,
	})
	cont = setupDependencies(cont)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"yarn", "plan"}})

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &BuildOutput{}, nil
}

type DeployInput struct {
	AwsAccessKeyIDSecretID     dagger.SecretID
	AwsSecretAccessKeySecretID dagger.SecretID
	AwsSessionTokenSecretID    dagger.SecretID
	CatApiKeySecretID          dagger.SecretID
	BackendDistDirectoryID     dagger.DirectoryID
}

type DeployOutput struct{}

func Deploy(ctx context.Context, client *dagger.Client, src dagger.DirectoryID, ipt *DeployInput) (*DeployOutput, error) {
	cont := newContainer(client, src, &newContainerInput{BackendDistDirectoryID: ipt.BackendDistDirectoryID})
	cont = setupSecrets(cont, &setupSecretsInput{
		AwsAccessKeyIDSecretID:     ipt.AwsAccessKeyIDSecretID,
		AwsSecretAccessKeySecretID: ipt.AwsSecretAccessKeySecretID,
		AwsSessionTokenSecretID:    ipt.AwsSessionTokenSecretID,
		CatApiKeySecretID:          ipt.CatApiKeySecretID,
	})
	cont = setupDependencies(cont)
	cont = cont.Exec(dagger.ContainerExecOpts{Args: []string{"yarn", "apply:auto-approve"}})

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}
	return &DeployOutput{}, nil
}

type newContainerInput struct {
	BackendDistDirectoryID dagger.DirectoryID
}

func newContainer(client *dagger.Client, src dagger.DirectoryID, ipt *newContainerInput) *dagger.Container {
	return util.
		NewContainer(client, src, "golang:1.19").
		WithEnvVariable("CI", "true").
		WithWorkdir("/app/infrastructure").
		WithMountedDirectory("/app/backend/dist", ipt.BackendDistDirectoryID)
}

type setupSecretsInput struct {
	AwsAccessKeyIDSecretID     dagger.SecretID
	AwsSecretAccessKeySecretID dagger.SecretID
	AwsSessionTokenSecretID    dagger.SecretID
	CatApiKeySecretID          dagger.SecretID
}

func setupSecrets(cont *dagger.Container, ipt *setupSecretsInput) *dagger.Container {
	return cont.
		WithSecretVariable("AWS_ACCESS_KEY_ID", ipt.AwsAccessKeyIDSecretID).
		WithSecretVariable("AWS_SECRET_ACCESS_KEY", ipt.AwsSecretAccessKeySecretID).
		WithSecretVariable("AWS_SESSION_TOKEN", ipt.AwsSessionTokenSecretID).
		WithSecretVariable("CAT_API_KEY", ipt.CatApiKeySecretID)
}

func setupDependencies(cont *dagger.Container) *dagger.Container {
	tfversion := "1.2.3"
	nodeversion := "16.x"

	cont = util.SetupUnzip(cont)
	cont = util.SetupTerraform(cont, tfversion)
	cont = util.SetupNodeJS(cont, nodeversion)
	cont = cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"go", "get"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"yarn", "install", "--frozen-lockfile"}})

	return cont
}
