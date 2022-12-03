package util

import (
	"fmt"

	"dagger.io/dagger"
)

func SetupNodeJS(cont *dagger.Container, version string) *dagger.Container {
	return cont.
		WithExec([]string{"bash", "-c", fmt.Sprintf("curl -fsSL https://deb.nodesource.com/setup_%s | bash -", version)}).
		WithExec([]string{"apt", "install", "-y", "nodejs"}).
		WithExec([]string{"npm", "install", "-g", "yarn"})
}
