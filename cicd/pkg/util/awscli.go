package util

import "dagger.io/dagger"

func SetupAWSCLI(cont *dagger.Container) *dagger.Container {
	return cont.
		Exec(dagger.ContainerExecOpts{Args: []string{
			"bash", "-c",
			`wget "https://awscli.amazonaws.com/awscli-exe-linux-$(uname -m).zip" -O /tmp/awscliv2.zip`,
		}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"unzip", "/tmp/awscliv2.zip", "-d", "/tmp"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"/tmp/aws/install"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"aws", "--version"}})
}
