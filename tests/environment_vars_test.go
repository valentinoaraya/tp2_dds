package tests

import (
	"testing"

	"github.com/valentinoaraya/tp2_dds/config"
)

func TestObtenerVariablesDeEntorno(t *testing.T) {

	vars := map[string]string{
		"DB_HOST":        config.DB_HOST,
		"DB_USER":        config.DB_USER,
		"DB_PASSWORD":    config.DB_PASSWORD,
		"DB_NAME":        config.DB_NAME,
		"DB_PORT":        config.DB_PORT,
		"CONTAINER_PORT": config.CONTAINER_PORT,
		"DB_SSL_MODE":    config.DB_SSL_MODE,
		"Url_connection": config.Url_connection,
	}

	for name, value := range vars {
		if value == "" {
			t.Errorf("La variable %s no se pudo cargar correctamente.", name)
		}
	}

}
