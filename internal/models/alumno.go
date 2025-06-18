package models

import "time"

type Alumno struct {
	Apellido        string
	Nombre          string
	NroDocumento    string
	TipoDocumento   string
	FechaNacimiento time.Time
	Sexo            string
	NroLegajo       string
	FechaIngreso    time.Time
}
