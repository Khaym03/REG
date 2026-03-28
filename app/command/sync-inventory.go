package command

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/domain"
)

type SyncInventoryCommand struct{}

type SyncInventoryHandler struct {
	scraper domain.InventoryScraper
	repo    domain.GuideRepository
}

func NewInventoryHandler(
	scraper domain.InventoryScraper,
	repo domain.GuideRepository,
) *SyncInventoryHandler {
	return &SyncInventoryHandler{
		scraper: scraper,
		repo:    repo,
	}
}

func (h *SyncInventoryHandler) Handle(ctx context.Context, cmd SyncInventoryCommand) error {
	remoteRubros, err := h.scraper.RubrosSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("snapshot inventory: %w", err)
	}

	remoteSet := make(map[string]struct{}, len(remoteRubros))
	for _, r := range remoteRubros {
		remoteSet[r.Name] = struct{}{}
	}

	localRubros := h.repo.GetRubros()

	// Insert missing rubros in remote
	for _, r := range localRubros {
		if _, exists := remoteSet[r.Name]; exists {
			continue
		}

		time.Sleep(time.Second * 5)
		if err := h.scraper.Insert(ctx, r); err != nil {
			return fmt.Errorf("insert rubro %s: %w", r.Name, err)
		}

		remoteSet[r.Name] = struct{}{}
	}

	return nil
}
