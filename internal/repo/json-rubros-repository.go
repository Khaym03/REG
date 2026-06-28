package repo

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"
)

var _ RubroRepository = (*JSONRubroRepository)(nil)

type JSONRubroRepository struct {
	store *JSONStore[RepositoryData]
}

func NewJSONRubroRepository(store *JSONStore[RepositoryData]) *JSONRubroRepository {
	return &JSONRubroRepository{store: store}
}

func (r *JSONRubroRepository) Save(
	ctx context.Context,
	rubros []domain.Rubro,
) error {
	return r.store.Update(func(data *RepositoryData) error {
		for _, rubro := range rubros {
			data.Rubros[rubro.Name] = rubro
		}

		return nil
	})
}

func (r *JSONRubroRepository) GetAll(
	ctx context.Context,
) (result []domain.Rubro, err error) {

	err = r.store.View(func(data RepositoryData) error {
		for _, r := range data.Rubros {
			result = append(result, r)
		}

		return nil
	})

	return
}
