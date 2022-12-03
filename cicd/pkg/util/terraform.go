package util

import (
	"fmt"

	"dagger.io/dagger"
)

func SetupTerraform(cont *dagger.Container, version string) *dagger.Container {
	return cont.
		WithExec([]string{
			"wget",
			fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_linux_amd64.zip", version, version),
			"-O",
			"/tmp/terraform.zip",
		}).
		WithExec([]string{"unzip", "/tmp/terraform.zip", "-d", "/usr/bin"}).
		WithExec([]string{"chmod", "+x", "/usr/bin/terraform"}).
		WithExec([]string{"terraform", "--version"})
}
