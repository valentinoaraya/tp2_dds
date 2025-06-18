package tests

import (
	"time"

	"github.com/valentinoaraya/tp2_dds/internal/models"
)

var InstanciaAlumno = models.Alumno{
	Apellido:        "Araya",
	Nombre:          "Valentino",
	NroDocumento:    "45361303",
	TipoDocumento:   "DNI",
	FechaNacimiento: time.Date(2004, 7, 14, 0, 0, 0, 0, time.UTC),
	Sexo:            "M",
	NroLegajo:       "9938",
	FechaIngreso:    time.Date(2022, 3, 8, 0, 0, 0, 0, time.UTC),
}

var Alumnos = []*models.Alumno{
	&InstanciaAlumno,
	{
		Apellido:        "Patiño",
		Nombre:          "Ignacio",
		NroDocumento:    "12345678",
		TipoDocumento:   "DNI",
		FechaNacimiento: time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC),
		Sexo:            "F",
		NroLegajo:       "1234",
		FechaIngreso:    time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		Apellido:        "Durán",
		Nombre:          "Faustino",
		NroDocumento:    "87654321",
		TipoDocumento:   "DNI",
		FechaNacimiento: time.Date(1999, 5, 20, 0, 0, 0, 0, time.UTC),
		Sexo:            "M",
		NroLegajo:       "5678",
		FechaIngreso:    time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		Apellido:        "Contreras",
		Nombre:          "Facundo",
		NroDocumento:    "1234123",
		TipoDocumento:   "DNI",
		FechaNacimiento: time.Date(1999, 5, 20, 0, 0, 0, 0, time.UTC),
		Sexo:            "M",
		NroLegajo:       "91011",
		FechaIngreso:    time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		Apellido:        "Romero",
		Nombre:          "Tomás",
		NroDocumento:    "12443474",
		TipoDocumento:   "DNI",
		FechaNacimiento: time.Date(1999, 5, 20, 0, 0, 0, 0, time.UTC),
		Sexo:            "M",
		NroLegajo:       "121314",
		FechaIngreso:    time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		Apellido:        "Perez",
		Nombre:          "Juan Ignacio",
		NroDocumento:    "09809789",
		TipoDocumento:   "DNI",
		FechaNacimiento: time.Date(1999, 5, 20, 0, 0, 0, 0, time.UTC),
		Sexo:            "M",
		NroLegajo:       "151617",
		FechaIngreso:    time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		Apellido:        "Vergara",
		Nombre:          "Juan",
		NroDocumento:    "6456345",
		TipoDocumento:   "DNI",
		FechaNacimiento: time.Date(1999, 5, 20, 0, 0, 0, 0, time.UTC),
		Sexo:            "M",
		NroLegajo:       "171819",
		FechaIngreso:    time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		Apellido:        "Campos",
		Nombre:          "Agustín",
		NroDocumento:    "543423",
		TipoDocumento:   "DNI",
		FechaNacimiento: time.Date(1999, 5, 20, 0, 0, 0, 0, time.UTC),
		Sexo:            "M",
		NroLegajo:       "202122",
		FechaIngreso:    time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
	},
}
