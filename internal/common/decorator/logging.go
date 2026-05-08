package decorator

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/common/decorator/command"
)

type LoggingDecorator[C any] struct {
	base command.CommandHandler[C]
}

func NewLoggingDecorator[C any](base command.CommandHandler[C]) LoggingDecorator[C] {
	return LoggingDecorator[C]{base: base}
}

func (d LoggingDecorator[C]) Handle(
	ctx context.Context,
	session auth.Session,
	cmd C,
) error {
	log.Infof("starting command %s\n", generateActionName(cmd))
	err := d.base.Handle(ctx, session, cmd)
	log.Infof("finished command %s\n", generateActionName(cmd))

	return err
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
