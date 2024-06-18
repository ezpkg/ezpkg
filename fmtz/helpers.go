package fmtz

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func mustZ(err error) {
	if err != nil {
		panic(err)
	}
}
