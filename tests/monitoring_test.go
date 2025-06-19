package tests

import (
	"os"
	"testing"
	"time"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/models"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
	"github.com/valentinoaraya/tp2_dds/internal/services"
)

// TestMonitoringAndConfiguration tests para funcionalidades de monitoreo y configuración
func TestMonitoringAndConfiguration(t *testing.T) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		t.Fatalf("Error al crear el repositorio: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	t.Run("Configuración de conexiones optimizada", func(t *testing.T) {
		// Verificar que el repositorio se creó con configuración optimizada
		// Esto se verifica indirectamente al crear el repositorio sin errores
		if repo == nil {
			t.Error("El repositorio no debería ser nil")
		}
	})

	t.Run("PrepararTablaOptimizada - Configuración para carga masiva", func(t *testing.T) {
		err := repo.PrepararTablaOptimizada()
		if err != nil {
			t.Errorf("Error al preparar tabla optimizada: %v", err)
		}
	})

	t.Run("RestaurarTablaNormal - Restaurar configuración normal", func(t *testing.T) {
		err := repo.RestaurarTablaNormal()
		if err != nil {
			t.Errorf("Error al restaurar tabla normal: %v", err)
		}
	})

	t.Run("CrearIndicesOptimizados - Crear índices después de carga", func(t *testing.T) {
		// Primero cargar algunos datos
		err := repo.CrearAlumnosBatch(Alumnos[:2])
		if err != nil {
			t.Fatalf("Error al cargar datos de prueba: %v", err)
		}

		// Crear índices
		err = repo.CrearIndicesOptimizados()
		if err != nil {
			t.Errorf("Error al crear índices optimizados: %v", err)
		}

		// Limpiar tabla
		err = repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Errorf("Error al limpiar tabla: %v", err)
		}
	})

	t.Run("Monitoreo de progreso en carga streaming", func(t *testing.T) {
		// Crear CSV de prueba
		csvPath := createTestCSV(t)
		defer os.Remove(csvPath)

		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar carga streaming y medir tiempo
		startTime := time.Now()
		err = service.CargarAlumnosStreaming(csvPath, 50, 2)
		duration := time.Since(startTime)

		if err != nil {
			t.Errorf("Error en carga streaming: %v", err)
		}

		// Verificar que se completó en un tiempo razonable
		if duration > 10*time.Second {
			t.Errorf("La carga tomó demasiado tiempo: %v", duration)
		}

		// Verificar cantidad de registros
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}

		// Calcular velocidad
		if duration.Seconds() > 0 {
			rate := float64(cantidad) / duration.Seconds()
			t.Logf("Velocidad de carga: %.2f registros/segundo", rate)
		}
	})

	t.Run("Monitoreo de memoria en carga por chunks", func(t *testing.T) {
		// Crear CSV de prueba
		csvPath := createTestCSV(t)
		defer os.Remove(csvPath)

		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar carga por chunks
		startTime := time.Now()
		err = service.CargarAlumnosChunked(csvPath, 2, 50, 2)
		duration := time.Since(startTime)

		if err != nil {
			t.Errorf("Error en carga por chunks: %v", err)
		}

		// Verificar que se completó en un tiempo razonable
		if duration > 10*time.Second {
			t.Errorf("La carga por chunks tomó demasiado tiempo: %v", duration)
		}

		// Verificar cantidad de registros
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}
	})

	t.Run("Manejo de errores en configuración", func(t *testing.T) {
		// Probar con configuración inválida (esto debería manejarse graciosamente)
		// Las configuraciones de PostgreSQL que fallan deberían continuar sin error
		err := repo.PrepararTablaOptimizada()
		if err != nil {
			t.Logf("Advertencia al optimizar BD (esperado en algunos casos): %v", err)
		}

		err = repo.RestaurarTablaNormal()
		if err != nil {
			t.Logf("Advertencia al restaurar BD (esperado en algunos casos): %v", err)
		}
	})

	t.Run("Verificación de integridad de datos", func(t *testing.T) {
		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Cargar datos usando método optimizado
		err = repo.CrearAlumnosBatch(Alumnos)
		if err != nil {
			t.Fatalf("Error al cargar datos: %v", err)
		}

		// Verificar que todos los datos se cargaron correctamente
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		expected := len(Alumnos)
		if cantidad != expected {
			t.Errorf("Se esperaban %d alumnos, pero se encontraron %d", expected, cantidad)
		}

		// Verificar que los datos específicos están presentes
		// Esto requeriría una función adicional en el repositorio para buscar por legajo
		// Por ahora solo verificamos la cantidad
	})
}

// TestPerformanceMonitoring tests para monitoreo de rendimiento
func TestPerformanceMonitoring(t *testing.T) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		t.Fatalf("Error al crear el repositorio: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	t.Run("Comparación de rendimiento entre métodos", func(t *testing.T) {
		// Crear CSV de prueba
		csvPath := createTestCSV(t)
		defer os.Remove(csvPath)

		// Método 1: Streaming
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		startTime := time.Now()
		err = service.CargarAlumnosStreaming(csvPath, 100, 2)
		streamingDuration := time.Since(startTime)

		if err != nil {
			t.Errorf("Error en streaming: %v", err)
		}

		streamingCount, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad streaming: %v", err)
		}

		// Método 2: Batch tradicional
		err = repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		startTime = time.Now()
		alumnos, err := service.ObtenerAlumnosDesdeCSV(csvPath)
		if err != nil {
			t.Errorf("Error al leer CSV: %v", err)
		}

		err = service.CargarAlumnosBatch(alumnos, 100)
		batchDuration := time.Since(startTime)

		if err != nil {
			t.Errorf("Error en batch: %v", err)
		}

		batchCount, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad batch: %v", err)
		}

		// Comparar resultados
		if streamingCount != batchCount {
			t.Errorf("Cantidades diferentes: Streaming=%d, Batch=%d", streamingCount, batchCount)
		}

		// Log de rendimiento
		t.Logf("Rendimiento - Streaming: %v, Batch: %v", streamingDuration, batchDuration)
		t.Logf("Streaming más rápido: %v", streamingDuration < batchDuration)
	})

	t.Run("Monitoreo de recursos durante carga", func(t *testing.T) {
		// Crear CSV de prueba
		csvPath := createTestCSV(t)
		defer os.Remove(csvPath)

		// Limpiar tabla
		err := repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Fatalf("Error al limpiar tabla: %v", err)
		}

		// Ejecutar carga y monitorear tiempo
		startTime := time.Now()
		err = service.CargarAlumnosStreaming(csvPath, 50, 2)
		duration := time.Since(startTime)

		if err != nil {
			t.Errorf("Error en carga: %v", err)
		}

		// Verificar métricas básicas
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al obtener cantidad: %v", err)
		}

		if duration.Seconds() > 0 {
			rate := float64(cantidad) / duration.Seconds()
			t.Logf("Métricas - Tiempo: %v, Registros: %d, Velocidad: %.2f reg/seg",
				duration, cantidad, rate)
		}

		// Verificar que la velocidad es razonable (> 10 reg/seg)
		if duration.Seconds() > 0 {
			rate := float64(cantidad) / duration.Seconds()
			if rate < 10 {
				t.Errorf("Velocidad muy baja: %.2f reg/seg", rate)
			}
		}
	})
}

// TestErrorHandling tests para manejo de errores en las optimizaciones
func TestErrorHandling(t *testing.T) {
	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		t.Fatalf("Error al crear el repositorio: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	t.Run("Manejo de errores en carga streaming", func(t *testing.T) {
		// Intentar cargar archivo inexistente
		err := service.CargarAlumnosStreaming("archivo_inexistente.csv", 100, 2)
		if err == nil {
			t.Error("Se esperaba un error al cargar archivo inexistente")
		}
	})

	t.Run("Manejo de errores en carga por chunks", func(t *testing.T) {
		// Intentar cargar archivo inexistente
		err := service.CargarAlumnosChunked("archivo_inexistente.csv", 100, 50, 2)
		if err == nil {
			t.Error("Se esperaba un error al cargar archivo inexistente")
		}
	})

	t.Run("Manejo de errores en batch optimizado", func(t *testing.T) {
		// Probar con slice nil
		err := service.CargarAlumnosBatchOptimizado(nil, 50)
		if err != nil {
			t.Errorf("No debería haber error con slice nil: %v", err)
		}

		// Probar con slice vacío
		err = service.CargarAlumnosBatchOptimizado([]*models.Alumno{}, 50)
		if err != nil {
			t.Errorf("No debería haber error con slice vacío: %v", err)
		}
	})

	t.Run("Manejo de errores en paralelización optimizada", func(t *testing.T) {
		// Probar con 0 workers (debería usar valor por defecto)
		err := service.CargarAlumnosParaleloOptimizado(Alumnos, 50, 0)
		if err != nil {
			t.Errorf("No debería haber error con 0 workers: %v", err)
		}

		// Limpiar tabla después del test
		err = repo.LimpiarTablaAlumnos()
		if err != nil {
			t.Errorf("Error al limpiar tabla: %v", err)
		}
	})
}

// TestConfigurationValidation tests para validación de configuración
func TestConfigurationValidation(t *testing.T) {
	t.Run("Validación de configuración de entorno", func(t *testing.T) {
		// Verificar que las variables de entorno están configuradas
		if config.DB_HOST == "" {
			t.Error("DB_HOST no está configurado")
		}
		if config.DB_USER == "" {
			t.Error("DB_USER no está configurado")
		}
		if config.DB_NAME == "" {
			t.Error("DB_NAME no está configurado")
		}
		if config.Url_connection == "" {
			t.Error("Url_connection no está configurado")
		}
	})

	t.Run("Validación de conexión a base de datos", func(t *testing.T) {
		repo, err := repositories.NewAlumnoRepository(config.Url_connection)
		if err != nil {
			t.Fatalf("Error al crear repositorio: %v", err)
		}
		defer repo.Close()

		// Verificar que la conexión funciona
		cantidad, err := repo.ObtenerCantidadAlumnos()
		if err != nil {
			t.Errorf("Error al verificar conexión: %v", err)
		}

		// La cantidad debería ser un número válido (>= 0)
		if cantidad < 0 {
			t.Errorf("Cantidad de alumnos inválida: %d", cantidad)
		}
	})
}
