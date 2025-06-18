package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	file, err := os.Create("data/alumnos_2.5M.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"apellido", "nombre", "nro_documento", "tipo_documento",
		"fecha_nacimiento", "sexo", "nro_legajo", "fecha_ingreso",
	})

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 2500000; i++ {
		apellido := fmt.Sprintf("Apellido%d", i)
		nombre := fmt.Sprintf("Nombre%d", i)
		nroDocumento := strconv.Itoa(30000000 + i)
		tipoDocumento := "DNI"
		fechaNacimiento := randomDate("1960-01-01", "2005-12-31", r)
		sexo := randomSexo(r)
		nroLegajo := strconv.Itoa(i)
		fechaIngreso := randomDate("2010-01-01", "2024-12-31", r)

		writer.Write([]string{
			apellido, nombre, nroDocumento, tipoDocumento,
			fechaNacimiento, sexo, nroLegajo, fechaIngreso,
		})

		if i%100000 == 0 {
			fmt.Printf("Generados %d alumnos...\n", i)
		}
	}

	fmt.Println("¡Archivo CSV generado con éxito!")
}

func randomDate(start, end string, r *rand.Rand) string {
	layout := "2006-01-02"
	startTime, _ := time.Parse(layout, start)
	endTime, _ := time.Parse(layout, end)

	delta := endTime.Sub(startTime).Hours() / 24
	randomDays := r.Intn(int(delta))
	randomDate := startTime.AddDate(0, 0, randomDays)

	return randomDate.Format(layout)
}

func randomSexo(r *rand.Rand) string {
	if r.Intn(2) == 0 {
		return "M"
	}
	return "F"
}
