package util

import (
	"fmt"

	"dagger.io/dagger"
)

func SetupTerraform(cont *dagger.Container, version string) *dagger.Container {
	return cont.
		Exec(dagger.ContainerExecOpts{Args: []string{
			"wget",
			fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_linux_amd64.zip", version, version),
			"-O",
			"/tmp/terraform.zip",
		}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"unzip", "/tmp/terraform.zip", "-d", "/usr/bin"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"chmod", "+x", "/usr/bin/terraform"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"terraform", "--version"}})
}
