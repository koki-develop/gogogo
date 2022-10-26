package util

import (
	"dagger.io/dagger"
)

func SetupTask(cont *dagger.Container) *dagger.Container {
	cont = cont.Exec(dagger.ContainerExecOpts{
		Args: []string{"go", "install", "github.com/go-task/task/v3/cmd/task@latest"},
	})
	cont = cont.Exec(dagger.ContainerExecOpts{
		Args: []string{"task", "--version"},
	})

	return cont
}
