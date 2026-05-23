package decorator

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/sirupsen/logrus"
)

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, session auth.Session, q Q) (R, error)
}

func ApplyQueryDecorators[H any, R any](
	handler QueryHandler[H, R],
	logger *logrus.Entry,
) QueryHandler[H, R] {

	return queryLoggingDecorator[H, R]{
		base:   handler,
		logger: logger,
	}
}
