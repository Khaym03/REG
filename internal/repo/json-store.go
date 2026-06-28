package repo

import (
	"sync"

	"github.com/Khaym03/REG/internal/domain"
)

type RepositoryData struct {
	Months         map[string][]domain.Guide         `json:"months"`
	Rubros         map[string]domain.Rubro           `json:"rubros"`
	ReceptionState map[string]domain.ReceptionResult `json:"reception_state"`
}

type JSONStore[T any] struct {
	mu sync.Mutex
	p  Persistence[T]
}

func NewJSONStore[T any](p Persistence[T]) *JSONStore[T] {
	return &JSONStore[T]{
		mu: sync.Mutex{},
		p:  p,
	}
}

func (s *JSONStore[T]) Update(fn func(*T) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.p.Load()
	if err != nil {
		return err
	}

	if err := fn(&data); err != nil {
		return err
	}

	return s.p.Save(data)
}

func (s *JSONStore[T]) View(fn func(T) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.p.Load()
	if err != nil {
		return err
	}

	return fn(data)
}
