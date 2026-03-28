package command

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/domain"
)

type SyncInventoryCommand struct{}

type SyncInventoryHandler struct {
	repo    domain.GuideRepository
	scraper domain.InventoryScraper
}

func NewInventoryHandler(
	repo domain.GuideRepository,
	scraper domain.InventoryScraper,
) *SyncInventoryHandler {
	return &SyncInventoryHandler{
		repo:    repo,
		scraper: scraper,
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
