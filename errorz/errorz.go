package errorz

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func MustZ(err error) {
	if err != nil {
		panic(err)
	}
}
