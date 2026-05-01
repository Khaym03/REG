package adapters

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Khaym03/REG/domain"
)

type repositoryData struct {
	Months         map[string][]domain.Guide         `json:"months"`
	Rubros         map[string]domain.Rubro           `json:"rubros"`
	ReceptionState map[string]domain.ReceptionResult `json:"reception_state"`
}

type JSONStore struct {
	filePath string
	mu       sync.Mutex
}

func NewJSONStore(filePath string) *JSONStore {
	return &JSONStore{filePath: filePath, mu: sync.Mutex{}}
}

func (s *JSONStore) Update(fn func(*repositoryData) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	if err := fn(&data); err != nil {
		return err
	}

	return s.save(data)
}

func (s *JSONStore) View(fn func(repositoryData) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	return fn(data)
}

func (s *JSONStore) save(data repositoryData) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}

func (s *JSONStore) load() (repositoryData, error) {
	var data repositoryData

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return repositoryData{
				Months:         make(map[string][]domain.Guide),
				Rubros:         make(map[string]domain.Rubro),
				ReceptionState: make(map[string]domain.ReceptionResult),
			}, nil
		}
		return data, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
