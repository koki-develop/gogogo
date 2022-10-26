package util

func Must[T any](rtn T, err error) T {
	if err != nil {
		panic(err)
	}

	return rtn
}
