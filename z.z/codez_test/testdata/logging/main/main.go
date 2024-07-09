package main

import (
	"fmt"

	"golang.org/x/net/context"

	"ezpkg.io/-/codez_test/testdata/logging"
)

var logger = logging.NewLogger(context.Background())

func main() {
	fmt.Println("Hello, World!")
	foo()
	bar()
}

func foo() {
	logger.Log("Hello %v!", "foo")
}

func bar() {
	l := logging.NewLogger(context.Background())
	l.Log("Goodbye %v!", "bar")

	l.LogX(context.Background(), "Hello %v!", "bar")
}
