package main

import (
	"fmt"
	"log"
	"time"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
	"github.com/valentinoaraya/tp2_dds/internal/services"
)

func main() {

	rutaArchivo := "data/alumnos_2.5M.csv"

	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	err = repo.LimpiarTablaAlumnos()
	if err != nil {
		log.Fatalf("Error limpiando tabla de alumnos: %v", err)
	}

	fmt.Println("Iniciando carga de alumnos desde CSV...")
	start := time.Now()
	service.CargarAlumnosStreaming(rutaArchivo, 5000, 16)
	tiempoFinal := time.Since(start)
	fmt.Printf("Carga de alumnos finalizada.\n")
	fmt.Printf("Tiempo total de carga: %v\n", tiempoFinal)
}
