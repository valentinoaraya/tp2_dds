# 🚀 TP2 - Carga Masiva de Alumnos

## 📋 Descripción

Este proyecto implementa un sistema de **carga masiva de datos de alumnos** en Go, diseñado para evaluar el rendimiento de inserción en base de datos PostgreSQL con archivos CSV grandes.

## 🏗️ Estructura del Proyecto

```
tp2_golang/
├── cmd/main.go                 # Punto de entrada de la aplicación
├── config/environment_vars.go  # Configuración de variables de entorno
├── data/alumnos.csv            # Datos de prueba (~10K registros)
├── data/alumnos_2.5M.csv       # Archivo grande (puede generarse localmente)
├── utils/generador_csv.go      # Generador de archivos CSV masivos
├── internal/models/alumno.go   # Modelo de datos Alumno
├── internal/repositories/      # Acceso a datos
├── internal/services/          # Lógica de negocio
├── init/init.sql               # Script de inicialización de BD
├── tests/                      # Tests unitarios
├── docker-compose.yml          # Configuración de Docker
└── go.mod                      # Dependencias de Go
```

## ⚙️ Configuración y Primeros Pasos

1. **Clonar el repositorio y ubicarse en la raíz.**
2. **Crear el archivo `.env` con los datos de conexión a la base de datos:**

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=tu_database
DB_PORT=5432
CONTAINER_PORT=5432
DB_SSL_MODE=disable
```

3. **Levantar la base de datos con Docker Compose:**

```bash
docker-compose up -d
```

4. **Instalar dependencias de Go:**

```bash
go mod download
```

5. **(Opcional) Generar el archivo grande de alumnos:**

```bash
go run utils/generador_csv.go
```
Esto crea `data/alumnos_2.5M.csv` (ya existe en el repo, pero puede regenerarse si se borra).

## 🏃‍♂️ Ejecución de la Aplicación

- **Por defecto, la aplicación procesa `data/alumnos_2.5M.csv` con batch de 5000 y 16 goroutines.**
- **Si querés usar el archivo chico, cambiá la variable `rutaArchivo` en `cmd/main.go` a `"data/alumnos.csv"`.**
- **Los parámetros de procesamiento (batch y goroutines) se modifican editando el código fuente.**

```bash
go run cmd/main.go
```

## 🧪 Tests

Para correr todos los tests:

```bash
go test ./...
```

## 🗄️ Estructura de la Tabla `Alumnos`

La tabla se crea automáticamente al iniciar el contenedor, con la siguiente estructura:

```sql
CREATE TABLE IF NOT EXISTS Alumnos (
    id SERIAL PRIMARY KEY,
    apellido VARCHAR(100) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    nro_documento VARCHAR(20) NOT NULL,
    tipo_documento VARCHAR(50) NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    sexo CHAR(1) NOT NULL,
    nro_legajo INTEGER NOT NULL,
    fecha_ingreso DATE NOT NULL
);
```

## 📝 Notas y Recomendaciones

- El archivo grande (`data/alumnos_2.5M.csv`) ocupa mucho espacio y memoria. Asegurate de tener recursos suficientes.
- El generador de CSV (`utils/generador_csv.go`) siempre genera 2.5 millones de registros y no acepta parámetros.
- No hay argumentos CLI ni flags: cualquier cambio de parámetros debe hacerse editando el código fuente (`cmd/main.go`).
- La conexión a la base de datos depende de las variables de entorno definidas en `.env`.

## 🐛 Troubleshooting

- **Error de conexión a BD:**
  - Verificá que Docker esté corriendo: `docker-compose ps`
  - Reiniciá el contenedor si es necesario: `docker-compose restart postgres`
- **Error de variables de entorno:**
  - Verificá el archivo `.env` y su ubicación en la raíz del proyecto.
- **Archivo CSV no encontrado:**
  - Generá el archivo con `go run utils/generador_csv.go` si falta.
- **Memoria insuficiente:**
  - Considerá aumentar la memoria disponible para Go o reducir el tamaño del archivo de entrada.

## 👨‍💻 Autor

**Valentino Araya** - Estudiante de Ingeniería en Sistemas de la Información - UTN FRSR

---

*Desarrollado con ❤️ en Go* 