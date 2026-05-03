package command

import (
	"context"

	"github.com/Khaym03/REG/domain"
)

type CommandHandler[C any] interface {
	Handle(ctx context.Context, session domain.Session, cmd C) error
}
