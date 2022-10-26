package backend

import (
	"context"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

type BuildOutput struct {
	Container       *dagger.Container
	DistDirectoryID dagger.DirectoryID
}

func Build(ctx context.Context, client *dagger.Client, src dagger.DirectoryID) (*BuildOutput, error) {
	// initialize
	cont := util.
		NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/backend")

	cont = util.SetupTask(cont)

	cont = cont.Exec(dagger.ContainerExecOpts{
		Args: []string{"task", "build"},
	})

	d, err := cont.Directory("/app/backend/dist").ID(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}

	return &BuildOutput{
		Container:       cont,
		DistDirectoryID: d,
	}, nil
}
