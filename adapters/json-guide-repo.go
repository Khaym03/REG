package adapters

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/utils"
)

var _ domain.GuideRepository = (*JSONGuideRepository)(nil)

const dateKeyFormat = "2006-01"

type repositoryData struct {
	Months         map[string][]domain.Guide         `json:"months"`
	Rubros         map[string]domain.Rubro           `json:"rubros"`
	ReceptionState map[string]domain.ReceptionResult `json:"reception_state"`
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

	key := date.From.Format(dateKeyFormat)
	_, exists := data.Months[key]

	return exists
}

func (r *JSONGuideRepository) SaveReceptionProgress(date utils.DateRange, result domain.ReceptionResult) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return
	}

	prev := data.ReceptionState[date.Key()]

	prev.Processed += result.Processed

	// if is already completed, do not broke it
	if result.Completed {
		prev.Completed = true
	}

	data.ReceptionState[date.Key()] = prev

	_ = r.save(data)
}

func (r *JSONGuideRepository) MarkReceptionCompleted(date utils.DateRange) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return
	}

	result := data.ReceptionState[date.Key()]
	result.Completed = true

	data.ReceptionState[date.Key()] = result

	_ = r.save(data)
}

func (r *JSONGuideRepository) IsReceptionCompleted(date utils.DateRange) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return false
	}

	result, exists := data.ReceptionState[date.Key()]
	if !exists {
		return false
	}

	return result.Completed
}

func (r *JSONGuideRepository) GetReceptionProgress(date utils.DateRange) domain.ReceptionResult {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.load()
	if err != nil {
		return domain.ReceptionResult{}
	}

	return data.ReceptionState[date.Key()]
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
