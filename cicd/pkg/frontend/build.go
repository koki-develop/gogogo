package frontend

import (
	"context"

	"dagger.io/dagger"
	"github.com/koki-develop/gogogo/cicd/pkg/util"
)

type BuildOutput struct{}

func Build(ctx context.Context, client *dagger.Client, src dagger.DirectoryID) (*BuildOutput, error) {
	// initialize
	cont := util.NewContainer(client, src, "golang:1.19").
		WithWorkdir("/app/frontend")

	// setup go-task
	cont = util.SetupTask(cont)

	// build
	cont = cont.Exec(dagger.ContainerExecOpts{
		Args: []string{"task", "build"},
	})

	// run
	if _, err := cont.ExitCode(ctx); err != nil {
		return nil, err
	}

	return &BuildOutput{}, nil
}
