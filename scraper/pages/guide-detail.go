package pages

import (
	"fmt"
	"strings"

	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

type GuideDetailsPage struct {
	page *rod.Page
}

func NewGuideDetailsPage(p *rod.Page) *GuideDetailsPage {
	return &GuideDetailsPage{page: p}
}

// ExtractRubros finds the "RUBROS" table and parses its content into domain entities
func (p *GuideDetailsPage) ExtractRubros() ([]domain.Rubro, error) {
	// Locate the header, then move to the table container
	// Using Search instead of MustElement to handle missing elements gracefully
	header, err := p.page.ElementR("h4", "RUBROS")
	if err != nil {
		return nil, fmt.Errorf("rubros header not found: %w", err)
	}

	// Navigate to the table sibling
	parent, err := header.Parent()
	if err != nil {
		return nil, err
	}
	grandParent, err := parent.Parent()
	if err != nil {
		return nil, err
	}
	tableContainer, err := grandParent.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to locate table container: %w", err)
	}

	rows, err := tableContainer.Elements("tbody tr")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch table rows: %w", err)
	}

	var results []domain.Rubro
	for i, row := range rows {
		rubro, err := p.parseRow(row)
		if err != nil {
			// We log and continue to avoid failing the entire guide for one bad row
			fmt.Printf("Warning: failed to parse row %d: %v\n", i, err)
			continue
		}
		results = append(results, rubro)
	}

	return results, nil
}

func (p *GuideDetailsPage) parseRow(row *rod.Element) (domain.Rubro, error) {
	cols, err := row.Elements("td")
	if err != nil {
		return domain.Rubro{}, err
	}

	if len(cols) < 1 {
		return domain.Rubro{}, fmt.Errorf("insufficient columns in row")
	}

	name, err := cols[0].Text()
	if err != nil {
		return domain.Rubro{}, fmt.Errorf("failed to get text from first column: %w", err)
	}

	return domain.Rubro{
		Name: strings.TrimSpace(name),
	}, nil
}
