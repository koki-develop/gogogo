package util

import (
	"fmt"

	"dagger.io/dagger"
)

func SetupNodeJS(cont *dagger.Container, version string) *dagger.Container {
	return cont.
		Exec(dagger.ContainerExecOpts{Args: []string{"bash", "-c", fmt.Sprintf("curl -fsSL https://deb.nodesource.com/setup_%s | bash -", version)}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"apt", "install", "-y", "nodejs"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"npm", "install", "-g", "yarn"}})
}
