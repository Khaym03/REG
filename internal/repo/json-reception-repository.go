package repo

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"
)

var _ ReceptionRepository = (*JSONReceptionRepository)(nil)

type JSONReceptionRepository struct {
	store *JSONStore
}

func NewJSONReceptionRepository(store *JSONStore) *JSONReceptionRepository {
	return &JSONReceptionRepository{store: store}
}

func (r *JSONReceptionRepository) SaveProgress(
	ctx context.Context,
	date domain.DateRange,
	result ReceptionResult,
) error {

	return r.store.Update(func(data *repositoryData) error {
		prev := data.ReceptionState[date.MonthKey()]

		prev.Processed += result.Processed

		// if is already completed, do not broke it
		if result.Completed {
			prev.Completed = true
		}

		data.ReceptionState[date.MonthKey()] = prev

		return nil
	})
}

func (r *JSONReceptionRepository) GetProgress(
	ctx context.Context,
	date domain.DateRange,
) (progress domain.ReceptionResult, err error) {

	err = r.store.View(func(data repositoryData) error {
		progress = data.ReceptionState[date.MonthKey()]

		return nil
	})

	return
}

func (r *JSONReceptionRepository) MarkCompleted(
	ctx context.Context,
	date domain.DateRange,
) error {

	return r.store.Update(func(data *repositoryData) error {
		result := data.ReceptionState[date.MonthKey()]
		result.Completed = true

		data.ReceptionState[date.MonthKey()] = result

		return nil
	})
}
func (r *JSONReceptionRepository) IsCompleted(
	ctx context.Context,
	date domain.DateRange,
) (completed bool, err error) {

	err = r.store.View(func(data repositoryData) error {
		result, exists := data.ReceptionState[date.MonthKey()]
		if !exists {
			completed = false
			return nil
		}

		completed = result.Completed

		return nil
	})

	return
}
