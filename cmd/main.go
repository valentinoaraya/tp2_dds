package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
	"github.com/valentinoaraya/tp2_dds/internal/services"
)

func main() {

	rutaArchivo := "data/alumnos_2.5M.csv"

	fmt.Println("🚀 Iniciando carga masiva de alumnos...")
	fmt.Println("=====================================")
	fmt.Printf("💻 CPU Cores disponibles: %d\n", runtime.NumCPU())
	fmt.Printf("🧠 Memoria disponible: %d MB\n", getMemoryInfo())
	fmt.Println()

	repo, err := repositories.NewAlumnoRepository(config.Url_connection)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer repo.Close()

	service := services.NewAlumnoService(repo)

	// Preparar base de datos para carga masiva
	fmt.Println("🔧 Preparando base de datos para carga masiva...")
	if err := repo.PrepararTablaOptimizada(); err != nil {
		log.Printf("⚠️  Advertencia al optimizar BD: %v", err)
	}

	// Estrategias optimizadas para 2.5M registros
	estrategias := []struct {
		nombre      string
		funcion     func() error
		descripcion string
	}{
		{
			nombre: "Streaming Ultra-Rápido (16 workers, 5000 batch)",
			funcion: func() error {
				return service.CargarAlumnosStreaming(rutaArchivo, 5000, 16)
			},
			descripcion: "Carga streaming con 16 workers y batch de 5000",
		},
		{
			nombre: "Streaming Conservador (8 workers, 2000 batch)",
			funcion: func() error {
				return service.CargarAlumnosStreaming(rutaArchivo, 2000, 8)
			},
			descripcion: "Carga streaming con 8 workers y batch de 2000",
		},
		{
			nombre: "Chunked Paralelo (100k chunks, 12 workers)",
			funcion: func() error {
				return service.CargarAlumnosChunked(rutaArchivo, 100000, 2000, 12)
			},
			descripcion: "Procesamiento por chunks de 100k con 12 workers",
		},
		{
			nombre: "Batch (1000 registros)",
			funcion: func() error {
				alumnos, err := service.ObtenerAlumnosDesdeCSV(rutaArchivo)
				if err != nil {
					return err
				}
				return service.CargarAlumnosBatch(alumnos, 1000)
			},
			descripcion: "Inserción por lotes de 1000 registros (método original)",
		},
		{
			nombre: "Paralelo (8 goroutines, 500 registros)",
			funcion: func() error {
				alumnos, err := service.ObtenerAlumnosDesdeCSV(rutaArchivo)
				if err != nil {
					return err
				}
				return service.CargarAlumnosParalelo(alumnos, 500, 8)
			},
			descripcion: "Inserción paralela con 8 goroutines (método original)",
		},
	}

	resultados := make([]struct {
		nombre      string
		tiempo      time.Duration
		descripcion string
		error       error
		registros   int
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
			registros   int
		}{
			nombre:      estrategia.nombre,
			tiempo:      duracion,
			descripcion: estrategia.descripcion,
			error:       err,
			registros:   cantidad,
		})

		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Completado en %v\n", duracion)
			fmt.Printf("   📊 Registros insertados: %d\n", cantidad)
			if duracion.Seconds() > 0 {
				rate := float64(cantidad) / duracion.Seconds()
				fmt.Printf("   ⚡ Velocidad: %.0f registros/segundo\n", rate)
			}
		}
		fmt.Println()
	}

	// Restaurar configuración normal de BD
	fmt.Println("🔧 Restaurando configuración normal de BD...")
	if err := repo.RestaurarTablaNormal(); err != nil {
		log.Printf("⚠️  Advertencia al restaurar BD: %v", err)
	}

	// Crear índices optimizados
	fmt.Println("🔧 Creando índices optimizados...")
	if err := repo.CrearIndicesOptimizados(); err != nil {
		log.Printf("⚠️  Error creando índices: %v", err)
	}

	fmt.Println("📈 RESUMEN DE RESULTADOS")
	fmt.Println("========================")
	fmt.Printf("%-50s %-15s %-15s %-10s\n", "Estrategia", "Tiempo", "Registros", "Estado")
	fmt.Println(strings.Repeat("-", 100))

	var mejorTiempo time.Duration
	var mejorEstrategia string
	var mejorRegistros int
	primerResultado := true

	for _, resultado := range resultados {
		estado := "✅ OK"
		if resultado.error != nil {
			estado = "❌ Error"
		} else if primerResultado || resultado.tiempo < mejorTiempo {
			mejorTiempo = resultado.tiempo
			mejorEstrategia = resultado.nombre
			mejorRegistros = resultado.registros
			primerResultado = false
		}

		fmt.Printf("%-50s %-15s %-15d %-10s\n",
			resultado.nombre,
			resultado.tiempo.String(),
			resultado.registros,
			estado,
		)
	}

	if !primerResultado {
		fmt.Println()
		fmt.Printf("🏆 Estrategia más rápida: %s\n", mejorEstrategia)
		fmt.Printf("⏱️  Tiempo: %v\n", mejorTiempo)
		fmt.Printf("📊 Registros: %d\n", mejorRegistros)
		if mejorTiempo.Seconds() > 0 {
			rate := float64(mejorRegistros) / mejorTiempo.Seconds()
			fmt.Printf("⚡ Velocidad: %.0f registros/segundo\n", rate)
		}
	}

	fmt.Println()
	fmt.Println("🎯 Recomendaciones para producción:")
	fmt.Println("   1. Usar la estrategia más rápida según resultados")
	fmt.Println("   2. Monitorear uso de memoria y CPU")
	fmt.Println("   3. Ajustar configuración de PostgreSQL según hardware")
	fmt.Println("   4. Considerar particionamiento para tablas > 10M registros")
	fmt.Println("   5. Implementar carga incremental para actualizaciones")
}

// getMemoryInfo obtiene información básica de memoria (simplificado)
func getMemoryInfo() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Sys / 1024 / 1024) // MB
}
