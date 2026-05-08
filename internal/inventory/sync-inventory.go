package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/internal/repo"
)

type SyncInventoryCommand struct{}

type SyncInventoryHandler struct {
	repo    repo.RubroRepository
	scraper InventoryService
}

func NewInventoryHandler(
	repo repo.RubroRepository,
	scraper InventoryService,
) *SyncInventoryHandler {
	return &SyncInventoryHandler{
		repo:    repo,
		scraper: scraper,
	}
}

func (h *SyncInventoryHandler) Handle(
	ctx context.Context,
	session Session,
	cmd SyncInventoryCommand,
) error {
	remoteRubros, err := h.scraper.Snapshot(ctx, session)
	if err != nil {
		return fmt.Errorf("snapshot inventory: %w", err)
	}

	remoteSet := make(map[string]struct{}, len(remoteRubros))
	for _, r := range remoteRubros {
		remoteSet[r.Name] = struct{}{}
	}

	localRubros, err := h.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Insert missing rubros in remote
	for _, r := range localRubros {
		if _, exists := remoteSet[r.Name]; exists {
			continue
		}

		time.Sleep(time.Second * 5)
		if err := h.scraper.Insert(ctx, session, r); err != nil {
			return fmt.Errorf("insert rubro %s: %w", r.Name, err)
		}

		remoteSet[r.Name] = struct{}{}
	}

	return nil
}
