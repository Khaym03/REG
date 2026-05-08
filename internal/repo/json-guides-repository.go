package repo

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"
)

var _ GuideRepository = (*JSONGuideRepository)(nil)

type JSONGuideRepository struct {
	store *JSONStore
}

func NewJSONGuideRepository(store *JSONStore) *JSONGuideRepository {
	return &JSONGuideRepository{store: store}
}

func (r *JSONGuideRepository) Exists(
	ctx context.Context,
	date domain.DateRange,
) (exists bool, err error) {

	err = r.store.View(func(data repositoryData) error {
		_, exists = data.Months[date.MonthKey()]
		return nil
	})

	return
}

func (r *JSONGuideRepository) Save(
	ctx context.Context,
	date domain.DateRange,
	guides []domain.Guide,
) error {

	return r.store.Update(func(data *repositoryData) error {
		data.Months[date.MonthKey()] = guides
		return nil
	})
}

func (r *JSONGuideRepository) Get(
	ctx context.Context,
	date domain.DateRange,
) (result []domain.Guide, err error) {

	err = r.store.View(func(data repositoryData) error {
		result = data.Months[date.MonthKey()]
		return nil
	})

	return
}
