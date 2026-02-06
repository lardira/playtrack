package middleware

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type TestHumaContext interface {
	huma.Context
}

type testCtx struct {
	TestHumaContext

	onSetStatus func(int)
	onHeader    func() string
}

func (t testCtx) Header(_ string) string {
	return t.onHeader()
}

func (t testCtx) SetStatus(code int) {
	t.onSetStatus(code)
}

func (t testCtx) Context() context.Context {
	return context.Background()
}
