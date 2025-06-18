package tests

import (
	"testing"

	"github.com/valentinoaraya/tp2_dds/config"
	"github.com/valentinoaraya/tp2_dds/internal/repositories"
)

func TestConexionDB(t *testing.T) {

	repo, err := repositories.NewAlumnoRepository(config.Url_connection)

	if err != nil {
		t.Errorf("Fallo al conectar a la base de datos: %v", err)
	}

	if repo == nil {
		t.Errorf("Repositorio es nil después de la conexión.")
	}

}
