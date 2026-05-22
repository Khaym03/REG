package decorator

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
)

type CommandHandler[C any] interface {
	Handle(ctx context.Context, session auth.Session, cmd C) error
}
