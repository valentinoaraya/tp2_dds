package tests

import (
	"os"
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

	t.Run("Cargar alumnos con Streaming de datos", func(t *testing.T) {
		// Crear un archivo CSV de prueba pequeño
		testCSVPath := "../data/test_alumnos.csv"
		testData := `apellido,nombre,nro_documento,tipo_documento,fecha_nacimiento,sexo,nro_legajo,fecha_ingreso
		Araya,Valentino,45361303,DNI,2004-07-14,M,9938,2022-03-08
		Patiño,Ignacio,12345678,DNI,2000-01-15,F,1234,2020-03-01
		Durán,Faustino,87654321,DNI,1999-05-20,M,5678,2021-08-15
		Contreras,Facundo,1234123,DNI,1999-05-20,M,91011,2021-08-15
		Romero,Tomás,12443474,DNI,1999-05-20,M,121314,2021-08-15
		Perez,Juan Ignacio,09809789,DNI,1999-05-20,M,151617,2021-08-15
		Vergara,Juan,6456345,DNI,1999-05-20,M,171819,2021-08-15
		Campos,Agustín,543423,DNI,1999-05-20,M,202122,2021-08-15`

		err := os.WriteFile(testCSVPath, []byte(testData), 0644)
		if err != nil {
			t.Fatalf("Error creando archivo CSV de prueba: %v", err)
		}
		defer os.Remove(testCSVPath)

		err = repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar la tabla: %v", err)
		}

		err = service.CargarAlumnosStreaming(testCSVPath, 3, 2)

		if err != nil {
			t.Errorf("Error al cargar alumnos: %v", err)
		}

		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad de alumnos: %v", err)
		}

		esperado := 8
		if cantidad != esperado {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", esperado, cantidad)
		}

		err = repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Errorf("Error al limpiar la tabla: %v", err)
		}

	})
}
