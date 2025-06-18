package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
	"github.com/valentinoaraya/tp2_dds/internal/services"
)

func main() {

	fmt.Println("🚀 Iniciando carga masiva de alumnos...")
	fmt.Println("=====================================")

	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	fmt.Println("🧹 Limpiando tabla de alumnos...")
	if err := repo.LimpiarTablaAlumnos(); err != nil {
		log.Fatalf("Error limpiando tabla: %v", err)
	}

	fmt.Println("📖 Cargando datos del archivo CSV...")
	alumnos, err := service.ObtenerAlumnosDesdeCSV("data/alumnos.csv")
	if err != nil {
		log.Fatalf("Error cargando CSV: %v", err)
	}

	fmt.Printf("✅ Cargados %d alumnos del CSV\n", len(alumnos))
	fmt.Println()

	estrategias := []struct {
		nombre      string
		funcion     func() error
		descripcion string
	}{
		{
			nombre: "Batch (100 registros)",
			funcion: func() error {
				return service.CargarAlumnosBatch(alumnos, 100)
			},
			descripcion: "Inserción por lotes de 100 registros",
		},
		{
			nombre: "Batch (500 registros)",
			funcion: func() error {
				return service.CargarAlumnosBatch(alumnos, 500)
			},
			descripcion: "Inserción por lotes de 500 registros",
		},
		{
			nombre: "Batch (1000 registros)",
			funcion: func() error {
				return service.CargarAlumnosBatch(alumnos, 1000)
			},
			descripcion: "Inserción por lotes de 1000 registros",
		},
		{
			nombre: "Paralelo (4 goroutines, 500 registros)",
			funcion: func() error {
				return service.CargarAlumnosParalelo(alumnos, 500, 4)
			},
			descripcion: "Inserción paralela con 4 goroutines",
		},
		{
			nombre: "Paralelo (8 goroutines, 500 registros)",
			funcion: func() error {
				return service.CargarAlumnosParalelo(alumnos, 500, 8)
			},
			descripcion: "Inserción paralela con 8 goroutines",
		},
	}

	resultados := make([]struct {
		nombre      string
		tiempo      time.Duration
		descripcion string
		error       error
	}, 0, len(estrategias))

	for _, estrategia := range estrategias {
		fmt.Printf("🔄 Ejecutando: %s\n", estrategia.nombre)
		fmt.Printf("   %s\n", estrategia.descripcion)

		if err := repo.LimpiarTablaAlumnos(); err != nil {
			log.Printf("Error limpiando tabla: %v", err)
			continue
		}

		inicio := time.Now()
		err := estrategia.funcion()
		duracion := time.Since(inicio)

		cantidad, countErr := repo.ObtenerCantidadAlumnos()
		if countErr != nil {
			log.Printf("Error contando registros: %v", countErr)
		}

		resultados = append(resultados, struct {
			nombre      string
			tiempo      time.Duration
			descripcion string
			error       error
		}{
			nombre:      estrategia.nombre,
			tiempo:      duracion,
			descripcion: estrategia.descripcion,
			error:       err,
		})

		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Completado en %v\n", duracion)
			fmt.Printf("   📊 Registros insertados: %d\n", cantidad)
		}
		fmt.Println()
	}

	fmt.Println("📈 RESUMEN DE RESULTADOS")
	fmt.Println("========================")
	fmt.Printf("%-40s %-15s %-10s\n", "Estrategia", "Tiempo", "Estado")
	fmt.Println(strings.Repeat("-", 70))

	for _, resultado := range resultados {
		estado := "✅ OK"
		if resultado.error != nil {
			estado = "❌ Error"
		}
		fmt.Printf("%-40s %-15s %-10s\n",
			resultado.nombre,
			resultado.tiempo.String(),
			estado)
	}

	var mejorTiempo time.Duration
	var mejorEstrategia string
	primerResultado := true

	for _, resultado := range resultados {
		if resultado.error == nil {
			if primerResultado || resultado.tiempo < mejorTiempo {
				mejorTiempo = resultado.tiempo
				mejorEstrategia = resultado.nombre
				primerResultado = false
			}
		}
	}

	if !primerResultado {
		fmt.Println()
		fmt.Printf("🏆 Estrategia más rápida: %s (%v)\n", mejorEstrategia, mejorTiempo)
		fmt.Printf("📊 Velocidad: %.2f registros/segundo\n", float64(len(alumnos))/mejorTiempo.Seconds())
	}

	fmt.Println()
	fmt.Println("🎯 Próximos pasos para escalar a 2.5M registros:")
	fmt.Println("   1. Ajustar tamaño de batch según resultados")
	fmt.Println("   2. Optimizar número de goroutines")
	fmt.Println("   3. Considerar particionamiento de tablas")
	fmt.Println("   4. Evaluar índices y configuración de PostgreSQL")
	fmt.Println("   5. Implementar carga incremental")
}
