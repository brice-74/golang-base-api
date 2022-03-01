package resolvers

import "context"

func (r Root) QueryCheck() string {
	return "ok"
}

func (r Root) QueryPanic(_ context.Context, params QueryPanicParams) string {
	if params.Panic {
		panic("I panic !!!")
	}

	return "No panic"
}

type QueryPanicParams struct {
	Panic bool
}
