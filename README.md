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
│   ├── alumnos.csv             # Datos de prueba (~10K registros)
│   └── alumnos_2.5M.csv        # Archivo grande (generado localmente)
├── utils/
│   └── generador_csv.go        # Generador de archivos CSV masivos
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

## 📁 Datos de Prueba

### **Archivos CSV Disponibles:**

1. **`data/alumnos.csv`** (10K registros)
   - ✅ **Incluido en el repositorio**
   - 📊 **10,001 registros** para pruebas rápidas
   - 🚀 **Ideal para desarrollo y testing**

2. **`data/alumnos_2.5M.csv`** (2.5M registros)
   - ⚠️ **NO incluido en el repositorio** (archivo muy pesado)
   - 📊 **2,500,000 registros** para pruebas de rendimiento
   - 🔧 **Debe generarse localmente** usando el generador

### **Generador de CSV Masivo**

Para generar el archivo de 2.5 millones de registros:

```bash
# Ejecutar el generador
go run utils/generador_csv.go
```

**Características del generador:**
- 🎲 **Datos aleatorios** pero realistas
- 📅 **Fechas de nacimiento** entre 1960-2005
- 📅 **Fechas de ingreso** entre 2010-2024
- 👥 **Distribución equilibrada** de géneros
- 📝 **Números de documento** únicos
- 🔢 **Números de legajo** secuenciales

**Progreso del generador:**
```
Generados 100000 alumnos...
Generados 200000 alumnos...
Generados 300000 alumnos...
...
Generados 2500000 alumnos...
¡Archivo CSV generado con éxito!
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

### **Preparación de Datos:**

```bash
# Para pruebas rápidas (10K registros)
# El archivo data/alumnos.csv ya está incluido

# Para pruebas de rendimiento (2.5M registros)
go run utils/generador_csv.go
```

### **Ejecutar la aplicación:**

```bash
# Con datos pequeños (10K registros)
go run cmd/main.go

# Con datos grandes (2.5M registros)
# Primero generar el archivo, luego ejecutar
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

### **Con 10K registros:**
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

### **Con 2.5M registros:**
```
🚀 Iniciando carga masiva de alumnos...
=====================================

🧹 Limpiando tabla de alumnos...
📖 Cargando datos del archivo CSV...
✅ Cargados 2500000 alumnos del CSV

🔄 Ejecutando: Batch (1000 registros)
   Inserción por lotes de 1000 registros
   ✅ Completado en 3m 45s
   📊 Registros insertados: 2500000

🏆 Estrategia más rápida: Paralelo (8 goroutines, 1000 registros) (2m 15s)
📊 Velocidad: 18518.52 registros/segundo
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

### **Archivo CSV no encontrado:**
```bash
# Verificar que existe el archivo
ls -la data/alumnos.csv

# Si no existe, generar el archivo grande
go run utils/generador_csv.go
```

### **Memoria insuficiente para 2.5M registros:**
```bash
# Aumentar memoria disponible para Go
export GOMEMLIMIT=4GiB

# O ejecutar con más memoria
go run -gcflags=-m cmd/main.go
```

## 📝 Licencia

Este proyecto es parte del Trabajo Práctico 2 de la materia Desarrollo de Software.

## 👨‍💻 Autor

**Valentino Araya** - Estudiante de Ingeniería en Sistemas de la Información - UTN FRSR

---

*Desarrollado con ❤️ en Go* 