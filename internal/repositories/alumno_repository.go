package repositories

import (
	"database/sql"
	"fmt"

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

	// Configuración optimizada para carga masiva
	db.SetMaxOpenConns(50)   // Más conexiones para paralelismo
	db.SetMaxIdleConns(25)   // Mantener conexiones activas
	db.SetConnMaxLifetime(0) // Sin límite de tiempo de vida

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

func (r *Repository) PrepararTablaOptimizada() error {
	queries := []string{
		"SET session_replication_role = replica;",
		"SET synchronous_commit = off;",
		"SET wal_buffers = '16MB';",
		"SET checkpoint_segments = 32;",
	}

	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			continue
		}
	}

	return nil
}

func (r *Repository) RestaurarTablaNormal() error {
	queries := []string{
		"SET session_replication_role = DEFAULT;",
		"SET synchronous_commit = on;",
	}

	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			continue
		}
	}

	return nil
}

func (r *Repository) CrearIndicesOptimizados() error {
	fmt.Println("🔧 Creando índices optimizados...")

	indices := []string{
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_alumnos_nro_legajo ON alumnos(nro_legajo);",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_alumnos_apellido ON alumnos(apellido);",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_alumnos_fecha_ingreso ON alumnos(fecha_ingreso);",
	}

	for _, index := range indices {
		if _, err := r.db.Exec(index); err != nil {
			fmt.Printf("⚠️  Error creando índice: %v\n", err)
			continue
		}
	}

	fmt.Println("✅ Índices creados")
	return nil
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
