package main

import (
	"context"
	"fmt"

	xcontext "golang.org/x/net/context"

	"ezpkg.io/-/codez_test/testpkg/logging"
)

var logger, _ = logging.NewLogger(context.Background())

func main() {
	fmt.Println("Hello, World!")
	foo()
	bar()

	logging.AliasCtx(context.Background())
	logging.GoOrgCtx(xcontext.Background())
}

func foo() error {
	logger.Log("Hello %v!", "foo")
	return nil
}

func bar() (string, error) {
	l, err := logging.NewLogger(context.Background())
	if err != nil {
		return "", err
	}
	l.Log("Goodbye %v!", "bar")

	l.LogX(context.Background(), "Hello %v!", "bar")
	return "bar", nil
}
