package decorator

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/sirupsen/logrus"
)

type CommandHandler[C any] interface {
	Handle(ctx context.Context, session auth.Session, cmd C) error
}

func ApplyCommandDecorators[H any](
	handler CommandHandler[H],
	logger *logrus.Entry,
) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base:   handler,
		logger: logger,
	}
}
