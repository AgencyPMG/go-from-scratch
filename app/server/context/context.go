package context

import (
	contextlib "golang.org/x/net/context"
)

type Context interface {
	contextlib.Context
}
