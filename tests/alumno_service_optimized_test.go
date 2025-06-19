package tests

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/models"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
	"github.com/valentinoaraya/tp2_dds/internal/services"
)

// TestAlumnoServiceOptimized tests para las nuevas funcionalidades optimizadas
func TestAlumnoServiceOptimized(t *testing.T) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		t.Fatalf("Error al crear el repositorio: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	// Crear archivo CSV de prueba para streaming
	csvPath := createTestCSV(t)
	defer os.Remove(csvPath)

	t.Run("CargarAlumnosStreaming - Carga streaming básica", func(t *testing.T) {
		// Limpiar tabla antes del test
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar carga streaming
		err = service.CargarAlumnosStreaming(csvPath, 100, 2)
		if err != nil {
			t.Errorf("Error en carga streaming: %v", err)
		}

		// Verificar que se cargaron los registros
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("CargarAlumnosStreaming - Con diferentes tamaños de batch", func(t *testing.T) {
		batchSizes := []int{50, 100, 200}

		for _, batchSize := range batchSizes {
			t.Run(fmt.Sprintf("Batch_%d", batchSize), func(t *testing.T) {
				// Limpiar tabla
				err := repo.LimpiarTablaAlumnos()
				if err != nil {
					t.Fatalf("Error al limpiar tabla: %v", err)
				}

				// Ejecutar carga
				err = service.CargarAlumnosStreaming(csvPath, batchSize, 2)
				if err != nil {
					t.Errorf("Error con batch %d: %v", batchSize, err)
				}

				// Verificar
				cantidad, err := repo.ObtenerCantidadAlumnos()
				if err != nil {
					t.Errorf("Error al obtener cantidad: %v", err)
				}

				expected := len(Alumnos)
				if cantidad != expected {
					t.Errorf("Batch %d: Se esperaban %d alumnos, pero se encontraron %d", batchSize, expected, cantidad)
				}
			})
		}
	})

	t.Run("CargarAlumnosStreaming - Con diferentes números de workers", func(t *testing.T) {
		workerCounts := []int{1, 2, 4}

		for _, workers := range workerCounts {
			t.Run(fmt.Sprintf("Workers_%d", workers), func(t *testing.T) {
				// Limpiar tabla
				err := repo.LimpiarTablaAlumnos()
				if err != nil {
					t.Fatalf("Error al limpiar tabla: %v", err)
				}

				// Ejecutar carga
				err = service.CargarAlumnosStreaming(csvPath, 100, workers)
				if err != nil {
					t.Errorf("Error con %d workers: %v", workers, err)
				}

				// Verificar
				cantidad, err := repo.ObtenerCantidadAlumnos()
				if err != nil {
					t.Errorf("Error al obtener cantidad: %v", err)
				}

				expected := len(Alumnos)
				if cantidad != expected {
					t.Errorf("Workers %d: Se esperaban %d alumnos, pero se encontraron %d", workers, expected, cantidad)
				}
			})
		}
	})

	t.Run("CargarAlumnosChunked - Carga por chunks", func(t *testing.T) {
		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar carga por chunks
		err = service.CargarAlumnosChunked(csvPath, 3, 50, 2)
		if err != nil {
			t.Errorf("Error en carga por chunks: %v", err)
		}

		// Verificar
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("CargarAlumnosParaleloOptimizado - Paralelización optimizada", func(t *testing.T) {
		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar carga paralela optimizada
		err = service.CargarAlumnosParaleloOptimizado(Alumnos, 50, 4)
		if err != nil {
			t.Errorf("Error en carga paralela optimizada: %v", err)
		}

		// Verificar
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("CargarAlumnosBatchOptimizado - Batch optimizado", func(t *testing.T) {
		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar batch optimizado
		err = service.CargarAlumnosBatchOptimizado(Alumnos, 50)
		if err != nil {
			t.Errorf("Error en batch optimizado: %v", err)
		}

		// Verificar
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("contarLineasCSV - Contar líneas del CSV", func(t *testing.T) {
		// En lugar de acceder al método privado, verificamos indirectamente
		// contando los registros cargados
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Cargar usando streaming
		err = service.CargarAlumnosStreaming(csvPath, 100, 2)
		if err != nil {
			t.Errorf("Error en carga streaming: %v", err)
		}

		// Verificar que se cargaron todos los registros
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("leerChunkCSV - Leer chunk específico", func(t *testing.T) {
		// En lugar de acceder al método privado, verificamos indirectamente
		// usando carga por chunks
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Cargar usando chunks pequeños
		err = service.CargarAlumnosChunked(csvPath, 2, 50, 2)
		if err != nil {
			t.Errorf("Error en carga por chunks: %v", err)
		}

		// Verificar que se cargaron todos los registros
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("CargarAlumnosStreaming - Manejo de errores", func(t *testing.T) {
		// Intentar cargar archivo inexistente
		err := service.CargarAlumnosStreaming("archivo_inexistente.csv", 100, 2)
		if err == nil {
			t.Error("Se esperaba un error al cargar archivo inexistente")
		}
	})

	t.Run("CargarAlumnosStreaming - CSV vacío", func(t *testing.T) {

		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Crear CSV vacío
		emptyCSV := createEmptyCSV(t)
		defer os.Remove(emptyCSV)

		err = service.CargarAlumnosStreaming(emptyCSV, 100, 2)
		if err != nil {
			t.Errorf("No debería haber error con CSV vacío: %v", err)
		}

		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		if cantidad != 0 {
			t.Errorf("Se esperaban 0 alumnos con CSV vacío, pero se encontraron %d", cantidad)
		}
	})

	t.Run("CargarAlumnosStreaming - CSV con registros corruptos", func(t *testing.T) {
		// Crear CSV con registros corruptos
		corruptCSV := createCorruptCSV(t)
		defer os.Remove(corruptCSV)

		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Debería continuar procesando a pesar de registros corruptos
		err = service.CargarAlumnosStreaming(corruptCSV, 100, 2)
		if err != nil {
			t.Errorf("No debería fallar con registros corruptos: %v", err)
		}

		// Debería haber cargado al menos algunos registros válidos
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		if cantidad == 0 {
			t.Error("Se esperaba que al menos algunos registros válidos se cargaran")
		}
	})
}

// TestRepositoryOptimized tests para las nuevas funcionalidades del repositorio
func TestRepositoryOptimized(t *testing.T) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		t.Fatalf("Error al crear el repositorio: %v", err)
	}
	defer repo.Close()

	t.Run("CrearAlumnosBatchOptimizado - Batch optimizado", func(t *testing.T) {
		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Crear batch optimizado
		err = repo.CrearAlumnosBatch(Alumnos)
		if err != nil {
			t.Errorf("Error en batch optimizado: %v", err)
		}

		// Verificar
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("PrepararTablaOptimizada - Configuración para carga masiva", func(t *testing.T) {
		err := repo.PrepararTablaOptimizada()
		if err != nil {
			t.Errorf("Error al preparar tabla optimizada: %v", err)
		}
	})

	t.Run("RestaurarTablaNormal - Restaurar configuración", func(t *testing.T) {
		err := repo.RestaurarTablaNormal()
		if err != nil {
			t.Errorf("Error al restaurar tabla normal: %v", err)
		}
	})

	t.Run("CrearIndicesOptimizados - Crear índices", func(t *testing.T) {
		err := repo.CrearIndicesOptimizados()
		if err != nil {
			t.Errorf("Error al crear índices optimizados: %v", err)
		}
	})

	t.Run("CrearAlumnosBatchOptimizado - Batch vacío", func(t *testing.T) {
		err := repo.CrearAlumnosBatch([]*models.Alumno{})
		if err != nil {
			t.Errorf("No debería haber error con batch vacío: %v", err)
		}
	})
}

// Funciones auxiliares para crear archivos CSV de prueba

func createTestCSV(t *testing.T) string {
	// Crear archivo CSV temporal
	tmpFile, err := os.CreateTemp("", "test_alumnos_*.csv")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	defer writer.Flush()

	// Escribir header
	header := []string{"apellido", "nombre", "nro_documento", "tipo_documento", "fecha_nacimiento", "sexo", "nro_legajo", "fecha_ingreso"}
	err = writer.Write(header)
	if err != nil {
		t.Fatalf("Error al escribir header: %v", err)
	}

	// Escribir datos de prueba
	for _, alumno := range Alumnos {
		row := []string{
			alumno.Apellido,
			alumno.Nombre,
			alumno.NroDocumento,
			alumno.TipoDocumento,
			alumno.FechaNacimiento.Format("2006-01-02"),
			alumno.Sexo,
			alumno.NroLegajo,
			alumno.FechaIngreso.Format("2006-01-02"),
		}
		err = writer.Write(row)
		if err != nil {
			t.Fatalf("Error al escribir fila: %v", err)
		}
	}

	return tmpFile.Name()
}

func createEmptyCSV(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "empty_alumnos_*.csv")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	defer writer.Flush()

	// Solo header
	header := []string{"apellido", "nombre", "nro_documento", "tipo_documento", "fecha_nacimiento", "sexo", "nro_legajo", "fecha_ingreso"}
	err = writer.Write(header)
	if err != nil {
		t.Fatalf("Error al escribir header: %v", err)
	}

	return tmpFile.Name()
}

func createCorruptCSV(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "corrupt_alumnos_*.csv")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	defer writer.Flush()

	// Header
	header := []string{"apellido", "nombre", "nro_documento", "tipo_documento", "fecha_nacimiento", "sexo", "nro_legajo", "fecha_ingreso"}
	err = writer.Write(header)
	if err != nil {
		t.Fatalf("Error al escribir header: %v", err)
	}

	// Algunos registros válidos
	for i, alumno := range Alumnos[:3] {
		row := []string{
			alumno.Apellido,
			alumno.Nombre,
			alumno.NroDocumento,
			alumno.TipoDocumento,
			alumno.FechaNacimiento.Format("2006-01-02"),
			alumno.Sexo,
			alumno.NroLegajo,
			alumno.FechaIngreso.Format("2006-01-02"),
		}
		err = writer.Write(row)
		if err != nil {
			t.Fatalf("Error al escribir fila válida %d: %v", i, err)
		}
	}

	// Registros corruptos
	corruptRows := [][]string{
		{"Solo", "un", "campo"},          // Muy pocos campos
		{"", "", "", "", "", "", "", ""}, // Campos vacíos
		{"Apellido", "Nombre", "123", "DNI", "fecha_invalida", "M", "1234", "2020-01-01"},      // Fecha inválida
		{"Apellido", "Nombre", "no_es_numero", "DNI", "2000-01-01", "M", "1234", "2020-01-01"}, // Número inválido
	}

	for _, row := range corruptRows {
		err = writer.Write(row)
		if err != nil {
			t.Fatalf("Error al escribir fila corrupta: %v", err)
		}
	}

	return tmpFile.Name()
}

// Benchmark tests para medir rendimiento

func BenchmarkCargarAlumnosStreaming(b *testing.B) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		b.Fatalf("Error al crear el repositorio: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)
	csvPath := createBenchmarkCSV(b)
	defer os.Remove(csvPath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			b.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar carga streaming
		err = service.CargarAlumnosStreaming(csvPath, 1000, 4)
		if err != nil {
			b.Fatalf("Error en benchmark: %v", err)
		}
	}
}

func BenchmarkCargarAlumnosBatchOptimizado(b *testing.B) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		b.Fatalf("Error al crear el repositorio: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			b.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar batch optimizado
		err = service.CargarAlumnosBatchOptimizado(Alumnos, 1000)
		if err != nil {
			b.Fatalf("Error en benchmark: %v", err)
		}
	}
}

func createBenchmarkCSV(b *testing.B) string {
	// Crear CSV más grande para benchmarks
	tmpFile, err := os.CreateTemp("", "benchmark_alumnos_*.csv")
	if err != nil {
		b.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	defer writer.Flush()

	// Header
	header := []string{"apellido", "nombre", "nro_documento", "tipo_documento", "fecha_nacimiento", "sexo", "nro_legajo", "fecha_ingreso"}
	err = writer.Write(header)
	if err != nil {
		b.Fatalf("Error al escribir header: %v", err)
	}

	// Crear 1000 registros para benchmark
	for i := 0; i < 1000; i++ {
		alumno := Alumnos[i%len(Alumnos)] // Reutilizar datos de prueba
		row := []string{
			fmt.Sprintf("%s_%d", alumno.Apellido, i),
			alumno.Nombre,
			fmt.Sprintf("%d", 10000000+i),
			alumno.TipoDocumento,
			alumno.FechaNacimiento.Format("2006-01-02"),
			alumno.Sexo,
			fmt.Sprintf("%d", 1000+i),
			alumno.FechaIngreso.Format("2006-01-02"),
		}
		err = writer.Write(row)
		if err != nil {
			b.Fatalf("Error al escribir fila %d: %v", i, err)
		}
	}

	return tmpFile.Name()
}
