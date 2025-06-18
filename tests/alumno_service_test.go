package tests

import (
	"testing"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
	"github.com/valentinoaraya/tp2_dds/internal/services"
)

func TestAlumnoService(t *testing.T) {

	repo, err := repositories.NewAlumnoRepository(config.Url_connection)

	if err != nil {
		t.Errorf("Error al crear el repositorio: %v", err)
	}

	service := services.NewAlumnoService(repo)

	t.Run("Crear instancia alumno", func(t *testing.T) {
		err = service.CrearAlumno(&InstanciaAlumno)

		if err != nil {
			t.Errorf("Error al crear el alumno: %v", err)
		}

		err = repo.LimpiarTablaAlumnos()

		if err != nil {
			t.Errorf("Error al limpiar la tabla alumnos: %v", err)
		}
	})

	t.Run("Obtener alumnos desde el CSV", func(t *testing.T) {
		alumnos, err := service.ObtenerAlumnosDesdeCSV("../data/alumnos.csv")

		if err != nil {
			t.Errorf("Error al obtener alumnos desde el CSV: %d", err)
		}

		if len(alumnos) <= 0 {
			t.Errorf("No se ha obtenido ningún alumno del archivo.")
		}

	})

	t.Run("Parsear un registro a un Alumno", func(t *testing.T) {

		registro := []string{
			"Araya",
			"Valentino",
			"45361303",
			"DNI",
			"2004-07-14",
			"M",
			"9938",
			"2022-03-08",
		}

		alumnoParseado, err := service.ParsearAlumno(registro)

		if err != nil {
			t.Errorf("Error al parsear el alumno: %v", err)
		}

		if alumnoParseado == nil {
			t.Error("El alumno parseado es 'nil.'")
		}
	})

	t.Run("Cargar un batch de alumnos", func(t *testing.T) {
		err := service.CargarAlumnosBatch(Alumnos, len(Alumnos))

		if err != nil {
			t.Errorf("Error al cargar batch de alumnos: %v", err)
		}

		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad de alumnos: %v", err)
		}

		esperado := len(Alumnos)
		if cantidad != esperado {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", esperado, cantidad)
		}

		err = repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Errorf("Error al limpiar la tabla: %v", err)
		}
	})

	t.Run("Cargar alumnos en paralelos", func(t *testing.T) {
		err := service.CargarAlumnosParalelo(Alumnos, 2, 2)

		if err != nil {
			t.Errorf("Ocurrió un error al cargar alumnos en paralelo: %v", err)
		}

		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad de alumnos: %v", err)
		}

		esperado := len(Alumnos)
		if cantidad != esperado {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", esperado, cantidad)
		}

		err = repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Errorf("Error al limpiar la tabla: %v", err)
		}
	})
}
