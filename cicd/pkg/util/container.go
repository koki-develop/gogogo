package util

import "dagger.io/dagger"

func NewContainer(client *dagger.Client, src dagger.DirectoryID, img string) *dagger.Container {
	return client.
		Container().
		From(img).
		WithMountedDirectory("/app", src)
}
