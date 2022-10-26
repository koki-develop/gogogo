package util

import (
	"dagger.io/dagger"
)

func SetupTask(cont *dagger.Container) *dagger.Container {
	return cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"go", "install", "github.com/go-task/task/v3/cmd/task@latest"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"task", "--version"}})
}
