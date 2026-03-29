package session

import (
	"context"
)

type key struct{}

func WithSession(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, key{}, s)
}

func FromContext(ctx context.Context) Session {
	s, ok := ctx.Value(key{}).(Session)
	if !ok {
		panic("session missing in context")
	}
	return s
}
