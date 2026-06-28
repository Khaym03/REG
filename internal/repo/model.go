package repo

import "github.com/Khaym03/REG/internal/domain"

type Persistence[T any] interface {
	Load() (T, error)
	Save(T) error
}

type (
	DateRange       = domain.DateRange
	ReceptionResult = domain.ReceptionResult
)
