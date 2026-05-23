package reception

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/constants"
	c "github.com/Khaym03/REG/internal/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

var _ Page = (*receptionPage)(nil)

type receptionPage struct {
	page *rod.Page
}

func NewReceptionPage(p *rod.Page) *receptionPage {
	return &receptionPage{page: p}
}

func (rp *receptionPage) Open() error {
	if err := browser.Navigate(rp.page, c.ReceptionURL); err != nil {
		return err
	}

	return rp.page.WaitIdle(c.TimeoutShort)
}

func (rp *receptionPage) ApplyFilters(date DateRange) (err error) {
	if err = browser.Click(rp.page, filterAccordionSelector); err != nil {
		return err
	}

	// Wait for the accordion animation to settle
	if err = rp.page.WaitDOMStable(c.TimeoutShort, 0.5); err != nil {
		return fmt.Errorf("page did not stabilize after opening filters: %w", err)
	}

	if err = browser.SelectOption(rp.page, selectStatusSelector, selectStatusOption); err != nil {
		return fmt.Errorf("failed to select status: %w", err)
	}
	if err = browser.SelectOption(rp.page, selectReceptionStatus, selectReceptionOption); err != nil {
		return fmt.Errorf("failed to select reception status: %w", err)
	}

	if err = browser.FillInputTime(rp.page, inputDateFromSelector, date.From); err != nil {
		return err
	}
	if err = browser.FillInputTime(rp.page, inputDateToSelector, date.To); err != nil {
		return err
	}

	if err = browser.Click(rp.page, filterButtonSelector); err != nil {
		return err
	}

	return rp.page.WaitDOMStable(c.TimeoutShort, 0.5)
}

func (rp *receptionPage) ConfirmReception() error {
	modal, err := rp.page.Timeout(constants.DefaultTimeout).Element(modalSelector)
	if err != nil {
		return err
	}
	// Wait for the modal to be fully visible
	if err := modal.Timeout(constants.DefaultTimeout).WaitVisible(); err != nil {
		return err
	}

	// Initialize a navigation waiter BEFORE the click.
	// This forces Rod to wait until the browser finishes the redirect/reload.
	wait := rp.page.MustWaitNavigation()

	log.Info("Confirming reception in modal...")
	if err := browser.Click(rp.page, modalConfirmBtnSelector); err != nil {
		return err
	}

	// Block execution here until navigation is complete
	wait()

	return nil
}

// Rows returns the collection of abstracted rows
func (rp *receptionPage) Rows() ([]TableRow, error) {
	elements, err := rp.page.Timeout(constants.DefaultTimeout).ElementsX(tableRowSelector)
	if !errors.Is(err, context.DeadlineExceeded) && err != nil {
		return nil, err
	}

	var rows []TableRow
	for _, el := range elements {
		rows = append(rows, &receptionRow{element: el})
	}

	log.Info("Rows: ", len(rows), rows)
	return rows, nil
}

var _ TableRow = (*receptionRow)(nil)

type receptionRow struct {
	element *rod.Element
}

func (r *receptionRow) ID() (string, error) {
	el, err := r.element.Timeout(constants.DefaultTimeout).ElementX(dataIDColumnSelector)
	if err != nil {
		return "", err
	}

	val, err := el.Timeout(constants.DefaultTimeout).Attribute("data-id_")
	if err != nil || val == nil {
		return "", err
	}
	return *val, nil
}

func (r *receptionRow) IsExpired() bool {
	badge, err := r.element.Timeout(constants.DefaultTimeout).Element(".badge-danger")
	if err != nil || badge == nil {
		return false
	}
	text, _ := badge.Text()
	return text == "VENCIDA"
}

func (r *receptionRow) TriggerReception() error {
	el, err := r.element.Timeout(constants.DefaultTimeout).Element(".recepcionar")
	if err != nil {
		return err
	}

	return el.Timeout(constants.DefaultTimeout).Click(proto.InputMouseButtonLeft, 1)
}

const (
	filterAccordionSelector = `//*[@id="accordion-filtros"]/div/div[1]/a`
	selectStatusSelector    = `//*[@id="select2-estatus_-container"]/..`
	selectReceptionStatus   = `//*[@id="select2-recepcion-container"]/..`
	selectStatusOption      = `//li[contains(text(), "APROBADA")]`
	selectReceptionOption   = `//li[contains(@id, "SIN_RECEPCIONAR")]`
	inputDateFromSelector   = `//*[@id="desde"]`
	inputDateToSelector     = `//*[@id="hasta"]`
	filterButtonSelector    = `//*[@id="collapse-filtro"]/div/form/div[3]/button`
	tableRowSelector        = `//table[@id="tabla-component"]/tbody/tr`
	dataIDColumnSelector    = `./td[@data-id_]`
	modalSelector           = "#recepcionarGuia"
	modalConfirmBtnSelector = "#recepcionarGuia button.btn-success"
)
