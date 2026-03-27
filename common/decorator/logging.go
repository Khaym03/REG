package decorator

import (
	"context"
	"log"

	"github.com/Khaym03/REG/common/decorator/command"
)

type LoggingDecorator[C any] struct {
	base command.CommandHandler[C]
}

func NewLoggingDecorator[C any](base command.CommandHandler[C]) LoggingDecorator[C] {
	return LoggingDecorator[C]{base: base}
}

func (d LoggingDecorator[C]) Handle(ctx context.Context, cmd C) error {
	log.Println("starting command")
	err := d.base.Handle(ctx, cmd)
	log.Println("finished command", err)
	return err
}
