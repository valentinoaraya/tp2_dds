package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valentinoaraya/tp2_dds/internal/models"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
)

type AlumnoService struct {
	repo *repositories.Repository
}

func NewAlumnoService(r *repositories.Repository) *AlumnoService {
	return &AlumnoService{repo: r}
}

func (s *AlumnoService) CrearAlumno(alumno *models.Alumno) error {

	if err := s.repo.CrearAlumno(alumno); err != nil {
		return fmt.Errorf("error creando alumno %s: %v", alumno.NroLegajo, err)
	}

	return nil
}

func (s *AlumnoService) ObtenerAlumnosDesdeCSV(rutaArchivo string) ([]*models.Alumno, error) {

	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo CSV: %v", err)
	}
	defer archivo.Close()

	lector := csv.NewReader(archivo)
	registros, err := lector.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error leyendo CSV: %v", err)
	}

	if len(registros) < 2 {
		return nil, fmt.Errorf("archivo CSV vacío o sin datos")
	}

	var alumnos []*models.Alumno
	for i, registro := range registros[1:] {
		if len(registro) < 8 {
			log.Printf("Registro %d incompleto, saltando: %v", i+1, registro)
			continue
		}

		alumno, err := s.ParsearAlumno(registro)
		if err != nil {
			log.Printf("Error parseando registro: %v", err)
			continue
		}

		alumnos = append(alumnos, alumno)
	}

	return alumnos, nil
}

func (s *AlumnoService) ParsearAlumno(registro []string) (*models.Alumno, error) {
	if len(registro) < 8 {
		return nil, fmt.Errorf("registro incompleto")
	}

	fechaNacimiento, err := time.Parse("2006-01-02", strings.TrimSpace(registro[4]))
	if err != nil {
		return nil, fmt.Errorf("error parseando fecha de nacimiento: %v", err)
	}

	fechaIngreso, err := time.Parse("2006-01-02", strings.TrimSpace(registro[7]))
	if err != nil {
		return nil, fmt.Errorf("error parseando fecha de ingreso: %v", err)
	}

	nroLegajo, err := strconv.Atoi(strings.TrimSpace(registro[6]))
	if err != nil {
		return nil, fmt.Errorf("error parseando número de legajo: %v", err)
	}

	return &models.Alumno{
		Apellido:        strings.TrimSpace(registro[0]),
		Nombre:          strings.TrimSpace(registro[1]),
		NroDocumento:    strings.TrimSpace(registro[2]),
		TipoDocumento:   strings.TrimSpace(registro[3]),
		FechaNacimiento: fechaNacimiento,
		Sexo:            strings.TrimSpace(registro[5]),
		NroLegajo:       strconv.Itoa(nroLegajo),
		FechaIngreso:    fechaIngreso,
	}, nil
}

func (s *AlumnoService) CargarAlumnosBatch(alumnos []*models.Alumno, tamanoBatch int) error {
	for i := 0; i < len(alumnos); i += tamanoBatch {
		fin := i + tamanoBatch
		if fin > len(alumnos) {
			fin = len(alumnos)
		}

		batch := alumnos[i:fin]
		if err := s.repo.CrearAlumnosBatch(batch); err != nil {
			return fmt.Errorf("error en batch %d-%d: %v", i, fin-1, err)
		}
	}
	return nil
}

func (s *AlumnoService) CargarAlumnosParalelo(alumnos []*models.Alumno, tamanoBatch int, numGoroutines int) error {
	if numGoroutines <= 0 {
		numGoroutines = 4
	}

	chunkSize := len(alumnos) / numGoroutines
	if chunkSize == 0 {
		chunkSize = 1
	}

	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if i == numGoroutines-1 {
			end = len(alumnos)
		}

		go func(start, end int) {
			defer wg.Done()
			chunk := alumnos[start:end]
			if err := s.CargarAlumnosBatch(chunk, tamanoBatch); err != nil {
				errChan <- fmt.Errorf("error en goroutine %d-%d: %v", start, end, err)
			}
		}(start, end)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// CargarAlumnosStreaming carga alumnos directamente desde CSV sin cargar todo en memoria
func (s *AlumnoService) CargarAlumnosStreaming(rutaArchivo string, tamanoBatch int, numGoroutines int) error {
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error abriendo archivo CSV: %v", err)
	}
	defer archivo.Close()

	// Contar líneas para el progreso
	totalLineas, err := s.contarLineasCSV(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error contando líneas: %v", err)
	}
	totalRegistros := totalLineas - 1 // Excluir header

	fmt.Printf("📊 Total de registros a procesar: %d\n", totalRegistros)

	// Configurar worker pool
	alumnosChan := make(chan []*models.Alumno, numGoroutines*2)
	errorChan := make(chan error, numGoroutines)
	var wg sync.WaitGroup

	// Iniciar workers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for batch := range alumnosChan {
				if err := s.repo.CrearAlumnosBatch(batch); err != nil {
					errorChan <- fmt.Errorf("worker %d error: %v", workerID, err)
					return
				}
			}
		}(i)
	}

	// Leer CSV y enviar batches
	lector := csv.NewReader(archivo)
	lector.FieldsPerRecord = -1 // Permitir campos variables

	// Saltar header
	if _, err := lector.Read(); err != nil {
		// Si el archivo está vacío (sin header), no es error grave
		return nil
	}

	var batch []*models.Alumno
	registrosProcesados := 0
	startTime := time.Now()
	alMenosUnValido := false

	for {
		registro, err := lector.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Error leyendo registro: %v", err)
			continue
		}

		alumno, err := s.ParsearAlumno(registro)
		if err != nil {
			log.Printf("Error parseando registro: %v", err)
			continue
		}

		alMenosUnValido = true
		batch = append(batch, alumno)
		registrosProcesados++

		if len(batch) >= tamanoBatch {
			select {
			case alumnosChan <- batch:
				// Mostrar progreso cada 100k registros
				if registrosProcesados%100000 == 0 {
					elapsed := time.Since(startTime)
					rate := float64(registrosProcesados) / elapsed.Seconds()
					fmt.Printf("📈 Progreso: %d/%d (%.1f%%) - %.0f reg/seg\n",
						registrosProcesados, totalRegistros,
						float64(registrosProcesados)/float64(totalRegistros)*100, rate)
				}
			case err := <-errorChan:
				close(alumnosChan)
				return err
			}
			batch = make([]*models.Alumno, 0, tamanoBatch)
		}
	}

	// Enviar batch final
	if len(batch) > 0 {
		alumnosChan <- batch
	}

	close(alumnosChan)
	wg.Wait()

	// Verificar errores de los workers
	select {
	case err := <-errorChan:
		return err
	default:
	}

	// Si no hubo ningún registro válido, no es error
	if !alMenosUnValido {
		return nil
	}

	totalTime := time.Since(startTime)
	fmt.Printf("✅ Carga completada en %v\n", totalTime)
	fmt.Printf("📊 Velocidad promedio: %.0f registros/segundo\n", float64(registrosProcesados)/totalTime.Seconds())

	return nil
}

// CargarAlumnosChunked carga el CSV en chunks para datasets muy grandes
func (s *AlumnoService) CargarAlumnosChunked(rutaArchivo string, tamanoChunk int, tamanoBatch int, numGoroutines int) error {
	totalLineas, err := s.contarLineasCSV(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error contando líneas: %v", err)
	}
	totalRegistros := totalLineas - 1

	fmt.Printf("📊 Total de registros: %d\n", totalRegistros)
	fmt.Printf("🔧 Configuración: Chunk=%d, Batch=%d, Workers=%d\n", tamanoChunk, tamanoBatch, numGoroutines)

	startTime := time.Now()
	registrosProcesados := 0

	// Procesar en chunks
	for offset := 0; offset < totalRegistros; offset += tamanoChunk {
		chunkEnd := offset + tamanoChunk
		if chunkEnd > totalRegistros {
			chunkEnd = totalRegistros
		}

		fmt.Printf("🔄 Procesando chunk %d-%d (%d registros)\n", offset, chunkEnd-1, chunkEnd-offset)

		alumnos, err := s.leerChunkCSV(rutaArchivo, offset, chunkEnd-offset)
		if err != nil {
			return fmt.Errorf("error leyendo chunk: %v", err)
		}

		if err := s.CargarAlumnosParaleloOptimizado(alumnos, tamanoBatch, numGoroutines); err != nil {
			return fmt.Errorf("error procesando chunk: %v", err)
		}

		registrosProcesados += len(alumnos)
		elapsed := time.Since(startTime)
		rate := float64(registrosProcesados) / elapsed.Seconds()
		fmt.Printf("📈 Progreso total: %d/%d (%.1f%%) - %.0f reg/seg\n",
			registrosProcesados, totalRegistros,
			float64(registrosProcesados)/float64(totalRegistros)*100, rate)
	}

	totalTime := time.Since(startTime)
	fmt.Printf("✅ Carga completada en %v\n", totalTime)
	fmt.Printf("📊 Velocidad promedio: %.0f registros/segundo\n", float64(registrosProcesados)/totalTime.Seconds())

	return nil
}

// CargarAlumnosParaleloOptimizado versión optimizada para grandes volúmenes
func (s *AlumnoService) CargarAlumnosParaleloOptimizado(alumnos []*models.Alumno, tamanoBatch int, numGoroutines int) error {
	if numGoroutines <= 0 {
		numGoroutines = 8
	}

	chunkSize := len(alumnos) / numGoroutines
	if chunkSize == 0 {
		chunkSize = 1
	}

	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if i == numGoroutines-1 {
			end = len(alumnos)
		}

		go func(start, end int) {
			defer wg.Done()
			chunk := alumnos[start:end]
			if err := s.CargarAlumnosBatchOptimizado(chunk, tamanoBatch); err != nil {
				errChan <- fmt.Errorf("error en goroutine %d-%d: %v", start, end, err)
			}
		}(start, end)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// CargarAlumnosBatchOptimizado usa el método optimizado del repositorio
func (s *AlumnoService) CargarAlumnosBatchOptimizado(alumnos []*models.Alumno, tamanoBatch int) error {
	for i := 0; i < len(alumnos); i += tamanoBatch {
		fin := i + tamanoBatch
		if fin > len(alumnos) {
			fin = len(alumnos)
		}

		batch := alumnos[i:fin]
		if err := s.repo.CrearAlumnosBatch(batch); err != nil {
			return fmt.Errorf("error en batch optimizado %d-%d: %v", i, fin-1, err)
		}
	}
	return nil
}

// contarLineasCSV cuenta las líneas del archivo CSV
func (s *AlumnoService) contarLineasCSV(rutaArchivo string) (int, error) {
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return 0, err
	}
	defer archivo.Close()

	lector := csv.NewReader(archivo)
	contador := 0

	// Leer línea por línea para ser más tolerante a errores
	for {
		_, err := lector.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			// Log del error pero continuar contando
			log.Printf("Error leyendo línea %d en conteo: %v", contador+1, err)
			contador++ // Contar la línea problemática también
			continue
		}
		contador++
	}

	return contador, nil
}

// leerChunkCSV lee un chunk específico del CSV
func (s *AlumnoService) leerChunkCSV(rutaArchivo string, offset, limit int) ([]*models.Alumno, error) {
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	lector := csv.NewReader(archivo)
	lector.FieldsPerRecord = -1

	// Saltar header
	if _, err := lector.Read(); err != nil {
		return nil, err
	}

	// Saltar hasta el offset
	for i := 0; i < offset; i++ {
		if _, err := lector.Read(); err != nil {
			return nil, err
		}
	}

	var alumnos []*models.Alumno
	for i := 0; i < limit; i++ {
		registro, err := lector.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Error leyendo registro en chunk: %v", err)
			continue
		}

		alumno, err := s.ParsearAlumno(registro)
		if err != nil {
			log.Printf("Error parseando registro en chunk: %v", err)
			continue
		}

		alumnos = append(alumnos, alumno)
	}

	return alumnos, nil
}
