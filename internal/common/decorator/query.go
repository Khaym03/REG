package decorator

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
)

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, session auth.Session, q Q) (R, error)
}
