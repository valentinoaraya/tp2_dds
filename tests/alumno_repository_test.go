package tests

import (
	"testing"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
)

func TestAlumnoRepository(t *testing.T) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)

	if err != nil {
		t.Errorf("Error creando el repository: %v", err)
	}

	t.Run("Limpiar tabla Alumnos", func(t *testing.T) {
		err = repo.LimpiarTablaAlumnos()

		if err != nil {
			t.Errorf("Error al limpiar la tabla: %v", err)
		}
	})

	t.Run("Persistir un alumno", func(t *testing.T) {

		err = repo.CrearAlumno(&InstanciaAlumno)

		if err != nil {
			t.Errorf("Error al crear el alumno: %v", err)
		}

		err = repo.LimpiarTablaAlumnos()

		if err != nil {
			t.Errorf("Error al limpiar la tabla: %v", err)
		}
	})

	t.Run("Inserta múltiples alumnos en un solo batch", func(t *testing.T) {
		err = repo.CrearAlumnosBatch(Alumnos)
		if err != nil {
			t.Errorf("Error al crear alumnos en batch: %v", err)
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
