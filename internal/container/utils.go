package container

import (
	"github.com/Khaym03/REG/internal/common/decorator"
	"github.com/Khaym03/REG/internal/common/decorator/command"
)

func withLogging[T any](handler command.CommandHandler[T]) command.CommandHandler[T] {
	return decorator.NewLoggingDecorator(handler)
}

func WithRetryAndLogging[T any](handler command.CommandHandler[T]) command.CommandHandler[T] {
	return decorator.NewLoggingDecorator(
		decorator.NewRetryDecorator(
			handler,
			decorator.DefaultRetryConfig,
		),
	)
}
