package receiver

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
)
