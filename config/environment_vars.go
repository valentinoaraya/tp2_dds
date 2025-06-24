package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	DB_HOST        string
	DB_USER        string
	DB_PASSWORD    string
	DB_NAME        string
	DB_PORT        string
	CONTAINER_PORT string
	DB_SSL_MODE    string
	Url_connection string
)

func init() {
	Init()
}

func Init() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error obteniendo directorio de trabajo: %v", err)
	}

	var envPath string
	if filepath.Base(workDir) == "tests" {
		envPath = filepath.Join(workDir, "..", ".env")
	} else {
		envPath = filepath.Join(workDir, ".env")
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error cargando el archivo .env desde %s: %v", envPath, err)
	}

	DB_HOST = os.Getenv("DB_HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	DB_PORT = os.Getenv("DB_PORT")
	CONTAINER_PORT = os.Getenv("CONTAINER_PORT")
	DB_SSL_MODE = os.Getenv("DB_SSL_MODE")
	Url_connection = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, CONTAINER_PORT, DB_SSL_MODE)
}
