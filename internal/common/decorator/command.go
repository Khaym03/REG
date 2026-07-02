package decorator

import (
	"context"

	"github.com/Khaym03/REG/internal/session"
	"github.com/sirupsen/logrus"
)

type CommandHandler[C any] interface {
	Handle(ctx context.Context, session Session, cmd C) error
}

type Session = session.Session

func ApplyCommandDecorators[H any](
	handler CommandHandler[H],
	logger *logrus.Entry,
) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base:   handler,
		logger: logger,
	}
}
