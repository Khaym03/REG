package pages

import (
	"fmt"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type ReceptionPage struct {
	page *rod.Page
}

func NewReceptionPage(p *rod.Page) *ReceptionPage {
	return &ReceptionPage{page: p}
}

func (rp *ReceptionPage) Open() error {
	if err := navigate(rp.page, c.ReceptionURL); err != nil {
		return err
	}

	return rp.page.WaitIdle(c.TimeoutShort)
}

func (rp *ReceptionPage) ApplyFilters(date domain.DateRange) (err error) {
	if err = click(rp.page, filterAccordionSelector); err != nil {
		return err
	}

	// Wait for the accordion animation to settle
	if err = rp.page.WaitDOMStable(c.TimeoutShort, 0.5); err != nil {
		return fmt.Errorf("page did not stabilize after opening filters: %w", err)
	}

	if err = selectOption(rp.page, selectStatusSelector, selectStatusOption); err != nil {
		return fmt.Errorf("failed to select status: %w", err)
	}
	if err = selectOption(rp.page, selectReceptionStatus, selectReceptionOption); err != nil {
		return fmt.Errorf("failed to select reception status: %w", err)
	}

	if err = fillInputTime(rp.page, inputDateFromSelector, date.From); err != nil {
		return err
	}
	if err = fillInputTime(rp.page, inputDateToSelector, date.To); err != nil {
		return err
	}

	if err = click(rp.page, filterButtonSelector); err != nil {
		return err
	}

	return rp.page.WaitDOMStable(c.TimeoutShort, 0.5)
}

func (rp *ReceptionPage) ConfirmReception() error {
	modal, err := rp.page.Element(modalSelector)
	if err != nil {
		return err
	}
	// Wait for the modal to be fully visible
	if err := modal.WaitVisible(); err != nil {
		return err
	}

	// Initialize a navigation waiter BEFORE the click.
	// This forces Rod to wait until the browser finishes the redirect/reload.
	wait := rp.page.MustWaitNavigation()

	fmt.Println("Confirming reception in modal...")
	if err := click(rp.page, modalConfirmBtnSelector); err != nil {
		return err
	}

	// Block execution here until navigation is complete
	wait()

	return nil
}

// Rows returns the collection of abstracted rows
func (rp *ReceptionPage) Rows() ([]*ReceptionRow, error) {
	elements, err := rp.page.ElementsX(tableRowSelector)
	if err != nil {
		return nil, err
	}

	var rows []*ReceptionRow
	for _, el := range elements {
		rows = append(rows, &ReceptionRow{element: el})
	}
	return rows, nil
}

// ReceptionRow represents a single line in the results table
type ReceptionRow struct {
	element *rod.Element
}

func (r *ReceptionRow) ID() string {
	val, _ := r.element.MustElement(dataIDColumnSelector).Attribute("data-id_")
	if val == nil {
		return ""
	}
	return *val
}

func (r *ReceptionRow) IsExpired() bool {
	badge, err := r.element.Element(".badge-danger")
	if err != nil || badge == nil {
		return false
	}
	text, _ := badge.Text()
	return text == "VENCIDA"
}

func (r *ReceptionRow) TriggerReception() error {
	el, err := r.element.Element(".recepcionar")
	if err != nil {
		return err
	}

	return el.Click(proto.InputMouseButtonLeft, 1)
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
