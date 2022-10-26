package util

import "dagger.io/dagger"

func SetupUnzip(cont *dagger.Container) *dagger.Container {
	return cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "update", "-qq"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "install", "-y", "unzip"}})
}
