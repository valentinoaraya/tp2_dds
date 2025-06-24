package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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

func (s *AlumnoService) ParsearAlumno(registro []string) (*models.Alumno, error) {
	if len(registro) < 8 {
		return nil, fmt.Errorf("registro incompleto")
	}

	fechaNacimiento, err := time.Parse("2006-01-02", registro[4])
	if err != nil {
		return nil, fmt.Errorf("error parseando fecha de nacimiento: %v", err)
	}

	fechaIngreso, err := time.Parse("2006-01-02", registro[7])
	if err != nil {
		return nil, fmt.Errorf("error parseando fecha de ingreso: %v", err)
	}

	nroLegajo, err := strconv.Atoi(registro[6])
	if err != nil {
		return nil, fmt.Errorf("error parseando número de legajo: %v", err)
	}

	return &models.Alumno{
		Apellido:        registro[0],
		Nombre:          registro[1],
		NroDocumento:    registro[2],
		TipoDocumento:   registro[3],
		FechaNacimiento: fechaNacimiento,
		Sexo:            registro[5],
		NroLegajo:       strconv.Itoa(nroLegajo),
		FechaIngreso:    fechaIngreso,
	}, nil
}

func (s *AlumnoService) CargarAlumnosStreaming(rutaArchivo string, tamanoBatch int, numGoroutines int) error {
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error abriendo archivo CSV: %v", err)
	}
	defer archivo.Close()

	alumnosChan := make(chan []*models.Alumno, numGoroutines*2)
	errorChan := make(chan error, numGoroutines)
	var wg sync.WaitGroup

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

	lector := csv.NewReader(archivo)
	lector.FieldsPerRecord = -1

	if _, err := lector.Read(); err != nil {
		return nil
	}

	var batch []*models.Alumno

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

		batch = append(batch, alumno)

		if len(batch) >= tamanoBatch {
			select {
			case err := <-errorChan:
				return err
			case alumnosChan <- batch:
				batch = make([]*models.Alumno, 0, tamanoBatch)
			}
		}
	}

	if len(batch) > 0 {
		select {
		case err := <-errorChan:
			return err
		case alumnosChan <- batch:
		}
	}

	close(alumnosChan)
	wg.Wait()

	select {
	case err := <-errorChan:
		return err
	default:
	}

	return nil
}
