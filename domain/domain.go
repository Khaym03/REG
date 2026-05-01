package domain

import (
	"errors"
	"fmt"

	"github.com/Khaym03/REG/constants"
)

type User struct {
	Username string
	Password string
}

type Rubro struct {
	Name string `json:"name"`
}

type Guide struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func NewGuide(id string) (Guide, error) {
	if id == "" {
		return Guide{}, errors.New("invalid ID for guide")
	}

	return Guide{
		ID:  id,
		URL: fmt.Sprintf("%s/%s", constants.GuidesURL, id),
	}, nil
}

type ReceptionResult struct {
	Processed int
	Completed bool
}
