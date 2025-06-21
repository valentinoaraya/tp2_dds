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

	rutaArchivo := "data/alumnos.csv"

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
			nombre: "Streaming de alumnos (16 workers, 5000 batch) Usando COPY",
			funcion: func() error {
				return service.CargarAlumnosStreaming(rutaArchivo, 5000, 16, "copy")
			},
		},
		{
			nombre: "Streaming de alumnos (16 workers, 5000 batch) Usando Un solo Insert",
			funcion: func() error {
				return service.CargarAlumnosStreaming(rutaArchivo, 5000, 16, "unInsert")
			},
		},
		{
			nombre: "Streaming de alumnos (16 workers, 5000 batch) Usando multiples Inserts",
			funcion: func() error {
				return service.CargarAlumnosStreaming(rutaArchivo, 5000, 16, "multiplesInserts")
			},
		},
	}

	resultados := make([]struct {
		nombre string
		tiempo time.Duration
		error  error
	}, 0, len(estrategias))

	for _, estrategia := range estrategias {
		fmt.Printf("🔄 Ejecutando: %s\n", estrategia.nombre)

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
			nombre string
			tiempo time.Duration
			error  error
		}{
			nombre: estrategia.nombre,
			tiempo: duracion,
			error:  err,
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
			estado,
		)
	}

	var mejorTiempo time.Duration
	var mejorEstrategia string
	var segundoMejorTiempo time.Duration
	var segundaMejorEstrategia string
	primerResultado := true
	segundoResultado := true

	for _, resultado := range resultados {
		if resultado.error == nil {
			if primerResultado || resultado.tiempo < mejorTiempo {
				if !primerResultado {
					segundoMejorTiempo = mejorTiempo
					segundaMejorEstrategia = mejorEstrategia
					segundoResultado = false
				}
				mejorTiempo = resultado.tiempo
				mejorEstrategia = resultado.nombre
				primerResultado = false
			} else if segundoResultado || resultado.tiempo < segundoMejorTiempo {
				segundoMejorTiempo = resultado.tiempo
				segundaMejorEstrategia = resultado.nombre
				segundoResultado = false
			}
		}
	}

	if !primerResultado {
		fmt.Println()
		fmt.Printf("🏆 Estrategia más rápida: %s (%v)\n", mejorEstrategia, mejorTiempo)
		fmt.Printf("📊 Velocidad: %.2f registros/segundo\n", float64(len(alumnos))/mejorTiempo.Seconds())
		fmt.Println()
	}

	if !segundoResultado {
		fmt.Printf("🥈 Segunda estrategia más rápida: %s (%v)\n", segundaMejorEstrategia, segundoMejorTiempo)
		fmt.Printf("📊 Velocidad: %.2f registros/segundo\n", float64(len(alumnos))/segundoMejorTiempo.Seconds())
		fmt.Println()
	}
}
