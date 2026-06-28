package repo

import (
	"encoding/json"
	"os"
)

type JSONPersistence[T any] struct {
	filePath string
	newValue func() T
}

func NewJSONPersistence[T any](filePath string, newValue func() T) *JSONPersistence[T] {
	return &JSONPersistence[T]{
		filePath: filePath,
		newValue: newValue,
	}
}

func (p *JSONPersistence[T]) Load() (T, error) {
	var data T

	file, err := os.Open(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return p.newValue(), nil
		}
		return data, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (p *JSONPersistence[T]) Save(data T) error {
	file, err := os.Create(p.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	return enc.Encode(data)
}
