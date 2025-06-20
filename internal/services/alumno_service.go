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
