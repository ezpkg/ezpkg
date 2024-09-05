package logging

import (
	stdctx "context"

	xctx "golang.org/x/net/context"
)

func AliasCtx(ctx stdctx.Context) error {
	return nil
}

func GoOrgCtx(ctx xctx.Context) error {
	return nil
}
