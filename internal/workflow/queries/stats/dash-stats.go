package stats

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/common/decorator"
	"github.com/Khaym03/REG/internal/event"
	"github.com/mustafaturan/bus/v3"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type StatsQuery struct{}

type StatsHandler decorator.QueryHandler[StatsQuery, Stats]

type statsHandler struct {
	logger   *logrus.Entry
	eventBus *bus.Bus
}

type Stats struct {
	OutstandingDebt   uint16 `json:"outstanding_debt"`
	InTransitGuides   uint16 `json:"intransit_guides"`
	ExpiredGuides     uint16 `json:"expired_guides"`
	PendingProcedures uint16 `json:"pending_procedures"`
}

func NewStatsHandler(logger *logrus.Entry, eventBus *bus.Bus) StatsHandler {
	return decorator.ApplyQueryDecorators(
		&statsHandler{
			logger:   logger,
			eventBus: eventBus,
		},
		logger,
	)
}

func (s Stats) IsZero() bool {
	return s.InTransitGuides == 0 && s.ExpiredGuides == 0
}

func (s Stats) String() string {
	var builder strings.Builder

	builder.Grow(120)

	fmt.Fprintf(&builder, "Outstanding Debt: %d\n", s.OutstandingDebt)
	fmt.Fprintf(&builder, "In-transit Guides: %d\n", s.InTransitGuides)
	fmt.Fprintf(&builder, "Expired Guides: %d\n", s.ExpiredGuides)
	fmt.Fprintf(&builder, "Pending Procedures: %d", s.PendingProcedures)

	return builder.String()
}

func (svc *statsHandler) Handle(
	ctx context.Context,
	session auth.Session,
	_ StatsQuery,
) (Stats, error) {
	var result Stats
	var mutex sync.Mutex

	session.Do(ctx, func(page *rod.Page) error {
		cards, err := page.Elements(cardSelector)
		if err != nil {
			return err
		}

		group, _ := errgroup.WithContext(ctx)

		for _, card := range cards {
			group.Go(func() error {
				statType, value, err := extractCardData(card)
				if err != nil {
					return err
				}

				mutex.Lock()
				result.add(statType, value)
				mutex.Unlock()

				return nil
			})
		}

		return group.Wait()
	})

	log.Info(result.String())

	svc.eventBus.Emit(ctx, event.StatsTopic, result)

	return result, nil
}

func extractCardData(card *rod.Element) (string, int, error) {
	title, err := extractCardTitle(card)
	if err != nil {
		return "", 0, err
	}

	value, err := extractCardValue(card, title)
	if err != nil {
		return "", 0, err
	}

	return title, value, nil
}

func extractCardTitle(card *rod.Element) (string, error) {
	element, err := card.Element("p")
	if err != nil {
		return "", err
	}

	text, err := element.Text()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

func extractCardValue(card *rod.Element, title string) (int, error) {
	element, err := card.Element("div")
	if err != nil {
		return 0, err
	}

	rawValue, err := element.Text()
	if err != nil {
		return 0, err
	}

	value := normalizeValue(rawValue, title)

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return parsedValue, nil
}

func normalizeValue(value, title string) string {
	value = strings.Split(value, "\n")[0]

	if title == labelOutstandingDebt {
		value = strings.ReplaceAll(value, "BS", "")
		value = strings.ReplaceAll(value, ",", "")
	}

	return strings.TrimSpace(value)
}

func (s Stats) HasActionableGuides(receiveInTransit bool) bool {
	hasExpiredGuides := s.ExpiredGuides > 0

	canReceiveInTransit := receiveInTransit &&
		s.InTransitGuides > 0

	return hasExpiredGuides || canReceiveInTransit
}

func (s *Stats) add(statType string, value int) {
	switch statType {
	case labelOutstandingDebt:
		s.OutstandingDebt += uint16(value)

	case labelInTransitGuides:
		s.InTransitGuides += uint16(value)

	case labelExpiredGuides:
		s.ExpiredGuides += uint16(value)

	case labelPendingProcedures:
		s.PendingProcedures += uint16(value)
	}
}

const cardSelector = `.card .text-white`

const (
	labelOutstandingDebt   = "DEUDA PENDIENTE"
	labelInTransitGuides   = "GUIAS EN TRANSITO"
	labelExpiredGuides     = "GUIAS VENCIDAS"
	labelPendingProcedures = "TRAMITES PENDIENTES"
)
