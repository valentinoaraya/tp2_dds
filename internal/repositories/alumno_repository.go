package repositories

import (
	"database/sql"

	"github.com/valentinoaraya/tp2_dds/internal/models"

	"github.com/lib/pq"
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

func (r *Repository) CrearAlumno(alumno *models.Alumno) error {
	query := `INSERT INTO Alumnos (apellido, nombre, nro_documento, tipo_documento, fecha_nacimiento, sexo, nro_legajo, fecha_ingreso) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query,
		alumno.Apellido,
		alumno.Nombre,
		alumno.NroDocumento,
		alumno.TipoDocumento,
		alumno.FechaNacimiento,
		alumno.Sexo,
		alumno.NroLegajo,
		alumno.FechaIngreso,
	)

	return err
}

func (r *Repository) CrearAlumnosBatch(alumnos []*models.Alumno) error {
	if len(alumnos) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(pq.CopyIn("alumnos", "apellido", "nombre", "nro_documento", "tipo_documento", "fecha_nacimiento", "sexo", "nro_legajo", "fecha_ingreso"))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, alumno := range alumnos {
		_, err := stmt.Exec(
			alumno.Apellido,
			alumno.Nombre,
			alumno.NroDocumento,
			alumno.TipoDocumento,
			alumno.FechaNacimiento.Format("2006-01-02"),
			alumno.Sexo,
			alumno.NroLegajo,
			alumno.FechaIngreso.Format("2006-01-02"),
		)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
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

func (r *Repository) Close() error {
	return r.db.Close()
}
