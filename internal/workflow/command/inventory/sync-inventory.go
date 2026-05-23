package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/internal/common/decorator"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/sirupsen/logrus"
)

type SyncInventoryCommand struct{}

type SyncInventoryHandler decorator.CommandHandler[SyncInventoryCommand]

type syncInventoryHandler struct {
	repo    repo.RubroRepository
	scraper InventoryService
}

func NewInventoryHandler(
	repo repo.RubroRepository,
	scraper InventoryService,
	logger *logrus.Entry,
) SyncInventoryHandler {
	return decorator.ApplyCommandDecorators(
		&syncInventoryHandler{
			repo:    repo,
			scraper: scraper,
		},
		logger,
	)
}

func (h *syncInventoryHandler) Handle(
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
