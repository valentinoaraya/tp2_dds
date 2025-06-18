# 🚀 TP2 - Carga Masiva de Alumnos

## 📋 Descripción

Este proyecto implementa un sistema de **carga masiva de datos de alumnos** en Go, diseñado para evaluar diferentes estrategias de inserción en base de datos PostgreSQL. El objetivo es optimizar el rendimiento para escalar a 2.5 millones de registros.

## 🎯 Objetivos

- **Evaluar estrategias de inserción masiva** en PostgreSQL
- **Comparar rendimiento** entre diferentes tamaños de batch
- **Analizar inserción paralela** con múltiples goroutines
- **Preparar el sistema** para escalar a 2.5M registros

## 🏗️ Arquitectura

```
tp2_golang/
├── cmd/
│   └── main.go                 # Punto de entrada de la aplicación
├── config/
│   └── environment_vars.go     # Configuración de variables de entorno
├── data/
│   └── alumnos.csv             # Datos de prueba (~10K registros)
├── internal/
│   ├── models/
│   │   └── alumno.go           # Modelo de datos Alumno
│   ├── repositories/
│   │   └── alumno_repository.go # Capa de acceso a datos
│   └── services/
│       └── alumno_service.go   # Lógica de negocio
├── init/
│   └── init.sql               # Script de inicialización de BD
├── tests/                     # Tests unitarios
├── docker-compose.yml         # Configuración de Docker
└── go.mod                     # Dependencias de Go
```

## 🚀 Estrategias Implementadas

### 1. **Inserción por Lotes (Batch)**
- **Batch 100**: Inserción en lotes de 100 registros
- **Batch 500**: Inserción en lotes de 500 registros  
- **Batch 1000**: Inserción en lotes de 1000 registros

### 2. **Inserción Paralela**
- **4 Goroutines**: Procesamiento paralelo con 4 workers
- **8 Goroutines**: Procesamiento paralelo con 8 workers

## 📊 Modelo de Datos

```go
type Alumno struct {
    Apellido        string
    Nombre          string
    NroDocumento    string
    TipoDocumento   string
    FechaNacimiento time.Time
    Sexo            string
    NroLegajo       string
    FechaIngreso    time.Time
}
```

## 🛠️ Tecnologías

- **Go 1.24.3**: Lenguaje principal
- **PostgreSQL 15**: Base de datos
- **Docker & Docker Compose**: Contenedores
- **godotenv**: Gestión de variables de entorno
- **lib/pq**: Driver de PostgreSQL para Go

## ⚙️ Configuración

### 1. **Variables de Entorno**

Crea un archivo `.env` en la raíz del proyecto:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=tu_database
DB_PORT=5432
CONTAINER_PORT=5432
DB_SSL_MODE=disable
```

### 2. **Inicialización con Docker**

```bash
# Levantar la base de datos
docker-compose up -d

# Verificar que esté corriendo
docker-compose ps
```

### 3. **Instalación de Dependencias**

```bash
go mod download
```

## 🏃‍♂️ Ejecución

### **Ejecutar la aplicación:**

```bash
go run cmd/main.go
```

### **Ejecutar tests:**

```bash
# Todos los tests
go test ./...

# Tests específicos
go test ./tests/...

# Con cobertura
go test -cover ./...
```

## 📈 Resultados Esperados

La aplicación ejecutará automáticamente todas las estrategias y mostrará:

```
🚀 Iniciando carga masiva de alumnos...
=====================================

🧹 Limpiando tabla de alumnos...
📖 Cargando datos del archivo CSV...
✅ Cargados 10001 alumnos del CSV

🔄 Ejecutando: Batch (100 registros)
   Inserción por lotes de 100 registros
   ✅ Completado en 2.5s
   📊 Registros insertados: 10001

📈 RESUMEN DE RESULTADOS
========================
Estrategia                              Tiempo         Estado
----------------------------------------------------------------------
Batch (100 registros)                   2.5s           ✅ OK
Batch (500 registros)                   1.8s           ✅ OK
Batch (1000 registros)                  1.2s           ✅ OK
Paralelo (4 goroutines, 500 registros)  0.9s           ✅ OK
Paralelo (8 goroutines, 500 registros)  0.7s           ✅ OK

🏆 Estrategia más rápida: Paralelo (8 goroutines, 500 registros) (0.7s)
📊 Velocidad: 14287.14 registros/segundo
```

## 🧪 Tests

El proyecto incluye tests unitarios para:

- ✅ **Modelos**: Validación de estructura de datos
- ✅ **Repositorios**: Operaciones de base de datos
- ✅ **Servicios**: Lógica de negocio
- ✅ **Configuración**: Variables de entorno
- ✅ **Utilidades**: Funciones auxiliares

### **Ejecutar tests específicos:**

```bash
# Tests de repositorio
go test ./tests/alumno_repository_test.go

# Tests de servicio
go test ./tests/alumno_service_test.go

# Tests de configuración
go test ./tests/environment_vars_test.go
```

## 🔧 Desarrollo

### **Estructura de Código**

- **Clean Architecture**: Separación clara de responsabilidades
- **Dependency Injection**: Inyección de dependencias
- **Error Handling**: Manejo robusto de errores
- **Logging**: Logs informativos para debugging

### **Patrones Utilizados**

- **Repository Pattern**: Abstracción de acceso a datos
- **Service Layer**: Lógica de negocio centralizada
- **Factory Pattern**: Creación de instancias
- **Strategy Pattern**: Diferentes estrategias de inserción

## 📊 Optimizaciones Implementadas

### **1. Inserción por Lotes**
```go
func (s *AlumnoService) CargarAlumnosBatch(alumnos []*models.Alumno, tamanoBatch int) error
```

### **2. Procesamiento Paralelo**
```go
func (s *AlumnoService) CargarAlumnosParalelo(alumnos []*models.Alumno, tamanoBatch int, numGoroutines int) error
```

### **3. Gestión de Recursos**
- Conexiones de base de datos optimizadas
- Uso eficiente de memoria
- Limpieza automática de recursos

## 🎯 Próximos Pasos para Escalar

1. **Ajustar tamaño de batch** según resultados obtenidos
2. **Optimizar número de goroutines** para el hardware específico
3. **Considerar particionamiento** de tablas
4. **Evaluar índices** y configuración de PostgreSQL
5. **Implementar carga incremental** para datasets grandes

## 🐛 Troubleshooting

### **Error de conexión a BD:**
```bash
# Verificar que Docker esté corriendo
docker-compose ps

# Reiniciar contenedor
docker-compose restart postgres
```

### **Error de variables de entorno:**
```bash
# Verificar archivo .env
cat .env

# Verificar que esté en la raíz del proyecto
ls -la .env
```

### **Error de permisos:**
```bash
# Dar permisos de ejecución
chmod +x cmd/main.go
```

## 📝 Licencia

Este proyecto es parte del Trabajo Práctico 2 de la materia Desarrollo de Software.

## 👨‍💻 Autor

**Valentino Araya** - Estudiante de Ingeniería en Sistemas de la Información - UTN FRSR

---

*Desarrollado con ❤️ en Go* 