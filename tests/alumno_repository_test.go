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

	t.Run("Recibe una lista de alumnos y los inserta con una sola consulta INSERT", func(t *testing.T) {
		err = repo.CrearAlumnosBatch(Alumnos)
		if err != nil {
			t.Errorf("Error al insertar alumnos con un solo INSERT: %v", err)
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
