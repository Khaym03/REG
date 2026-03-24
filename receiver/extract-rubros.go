package receiver

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
)

type Rubro struct {
	Nombre       string
	Cantidad     string
	PrecioVenta  string
	Presentacion string
	Marca        string
}

func ExtractRubrosFromGuide(page *rod.Page) []Rubro {
	table := page.MustElementR("h4", "RUBROS").MustParent().MustParent().MustNext()

	rows := table.MustElements("tbody tr")

	var listaRubros []Rubro

	for _, row := range rows {
		cols := row.MustElements("td")

		// Ensure the row has the expected columns
		if len(cols) >= 5 {
			item := Rubro{
				Nombre:       strings.TrimSpace(cols[0].MustText()),
				Cantidad:     strings.TrimSpace(cols[1].MustText()),
				PrecioVenta:  strings.TrimSpace(cols[2].MustText()),
				Presentacion: strings.TrimSpace(cols[3].MustText()),
				Marca:        strings.TrimSpace(cols[4].MustText()),
			}
			listaRubros = append(listaRubros, item)
		}
	}

	for _, r := range listaRubros {
		fmt.Printf("Producto: %s | Cantidad: %s\n", r.Nombre, r.Cantidad)
	}

	return listaRubros
}
