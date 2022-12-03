package util

import "dagger.io/dagger"

func SetupUnzip(cont *dagger.Container) *dagger.Container {
	return cont.
		WithExec([]string{"apt", "update", "-qq"}).
		WithExec([]string{"apt", "install", "-y", "unzip"})
}
