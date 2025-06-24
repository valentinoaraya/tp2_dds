package repositories

import (
	"fmt"
	"strings"

	"database/sql"

	"github.com/valentinoaraya/tp2_dds/internal/models"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewAlumnoRepository(connStr string) (*Repository, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &Repository{db: db}, nil
}

func (r *Repository) LimpiarTablaAlumnos() error {
	query := "TRUNCATE TABLE alumnos"

	_, err := r.db.Exec(query)

	return err
}

func (r *Repository) ObtenerCantidadAlumnos() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM Alumnos").Scan(&count)
	return count, err
}

func (r *Repository) CrearAlumnosBatch(alumnos []*models.Alumno) error {
	if len(alumnos) == 0 {
		return nil
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO alumnos (")
	queryBuilder.WriteString("apellido, nombre, nro_documento, tipo_documento, fecha_nacimiento, sexo, nro_legajo, fecha_ingreso")
	queryBuilder.WriteString(") VALUES ")
	valores := []interface{}{}

	for i, alumno := range alumnos {
		n := i * 8
		queryBuilder.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8))

		valores = append(valores,
			alumno.Apellido,
			alumno.Nombre,
			alumno.NroDocumento,
			alumno.TipoDocumento,
			alumno.FechaNacimiento,
			alumno.Sexo,
			alumno.NroLegajo,
			alumno.FechaIngreso,
		)
	}

	query := strings.TrimSuffix(queryBuilder.String(), ",")

	_, err := r.db.Exec(query, valores...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}
