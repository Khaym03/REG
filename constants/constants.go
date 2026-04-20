package constants

import "time"

const (
	BaseURL      = `https://sica.sunagro.gob.ve`
	ReceptionURL = BaseURL + `/despachos/recepcion`
	InventoryURL = BaseURL + `/inventario`
	GuidesURL    = BaseURL + `/guias`
	LoginURL     = BaseURL + `/login`

	DateKeyFormat = "2006-01"

	TimeoutMedium = time.Second * 5
)
