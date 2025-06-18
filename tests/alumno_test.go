package tests

import (
	"testing"
)

type test struct {
	obtenido string
	esperado string
}

func TestAlumno(t *testing.T) {

	tests := map[string]test{
		"Validar Apellido":            {obtenido: InstanciaAlumno.Apellido, esperado: "Araya"},
		"Validar Nombre":              {obtenido: InstanciaAlumno.Nombre, esperado: "Valentino"},
		"Validar Número de documento": {obtenido: InstanciaAlumno.NroDocumento, esperado: "45361303"},
		"Validar Tipo de documento":   {obtenido: InstanciaAlumno.TipoDocumento, esperado: "DNI"},
		"Validar Fecha de nacimeinto": {obtenido: InstanciaAlumno.FechaNacimiento.Format("2006-01-02"), esperado: "2004-07-14"},
		"Validar Sexo":                {obtenido: InstanciaAlumno.Sexo, esperado: "M"},
		"Validar Número de legajo":    {obtenido: InstanciaAlumno.NroLegajo, esperado: "9938"},
		"Validar Fecha de ingreso":    {obtenido: InstanciaAlumno.FechaIngreso.Format("2006-01-02"), esperado: "2022-03-08"},
	}

	for nombre, test := range tests {
		t.Run(nombre, func(t *testing.T) {
			if test.esperado != test.obtenido {
				t.Errorf("Se esperaba que '%s' sea '%s', pero se obtuvo '%s'", nombre, test.esperado, test.obtenido)
			}
		})
	}

}
