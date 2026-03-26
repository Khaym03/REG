package decorator

import (
	"context"
	"log"
)

type LoggingDecorator[C any] struct {
	base CommandHandler[C]
}

func (d LoggingDecorator[C]) Handle(ctx context.Context, cmd C) error {
	log.Println("starting command")
	err := d.base.Handle(ctx, cmd)
	log.Println("finished command", err)
	return err
}
