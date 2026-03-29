package decorator

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Khaym03/REG/common/decorator/command"
)

type LoggingDecorator[C any] struct {
	base command.CommandHandler[C]
}

func NewLoggingDecorator[C any](base command.CommandHandler[C]) LoggingDecorator[C] {
	return LoggingDecorator[C]{base: base}
}

func (d LoggingDecorator[C]) Handle(ctx context.Context, cmd C) error {
	log.Printf("starting command %s\n", generateActionName(cmd))
	err := d.base.Handle(ctx, cmd)
	log.Printf("finished command %s\n", generateActionName(cmd))

	return err
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
