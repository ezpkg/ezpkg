package decorator

import "context"

// +ggen:decorator
type MyInterface interface {
	GetInt() int

	GetString() string

	DoSomething(ctx context.Context, a int, b string, s MyStruct) (MyStruct, error)
}

type MyStruct struct {
}
