package adapters

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/utils"
)

var _ domain.GuideRepository = (*JSONGuideRepository)(nil)

type repositoryData struct {
	Months map[string][]domain.Guide `json:"months"`
	Rubros map[string]domain.Rubro   `json:"rubros"`
}

type JSONGuideRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewJSONGuideRepository(filePath string) *JSONGuideRepository {
	return &JSONGuideRepository{filePath: filePath, mu: sync.Mutex{}}
}

func (r *JSONGuideRepository) Exists(date utils.DateRange) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return false
	}

	key := date.From.Format("2006-01")
	_, exists := data.Months[key]

	return exists
}

func (r *JSONGuideRepository) SaveGuides(date utils.DateRange, guides []domain.Guide) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return
	}

	key := date.From.Format("2006-01")
	data.Months[key] = guides

	_ = r.save(data)
}

func (r *JSONGuideRepository) GetGuides(date utils.DateRange) []domain.Guide {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return nil
	}

	key := date.From.Format("2006-01")
	return data.Months[key]
}

func (r *JSONGuideRepository) SaveRubros(rubros []domain.Rubro) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return
	}

	for _, r := range rubros {
		data.Rubros[r.Name] = r
	}

	_ = r.save(data)
}

func (r *JSONGuideRepository) GetRubros() []domain.Rubro {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return nil
	}

	var result []domain.Rubro
	for _, r := range data.Rubros {
		result = append(result, r)
	}

	return result
}

func (r *JSONGuideRepository) load() (repositoryData, error) {
	var data repositoryData

	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return repositoryData{
				Months: make(map[string][]domain.Guide),
				Rubros: make(map[string]domain.Rubro),
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

func (r *JSONGuideRepository) save(data repositoryData) error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}
