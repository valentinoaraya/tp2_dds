# Makefile para TP2 Golang - Carga Masiva de Alumnos

.PHONY: build run test clean docker-up docker-down help monitor test-optimized test-all

# Variables
BINARY_NAME=alumnos-loader
MAIN_FILE=cmd/main.go
DOCKER_COMPOSE_FILE=docker-compose.yml

# Compilar el programa
build:
	@echo "🔨 Compilando programa..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Compilación completada: $(BINARY_NAME)"

# Ejecutar el programa
run: build
	@echo "🚀 Ejecutando programa..."
	./$(BINARY_NAME)

# Ejecutar sin compilar (si ya existe el binario)
run-only:
	@echo "🚀 Ejecutando programa..."
	./$(BINARY_NAME)

# Ejecutar con datos de prueba (CSV pequeño)
run-test:
	@echo "🧪 Ejecutando con datos de prueba..."
	go run $(MAIN_FILE)

# Ejecutar tests básicos
test:
	@echo "🧪 Ejecutando tests básicos..."
	go test ./tests -v

# Ejecutar tests de optimizaciones
test-optimized:
	@echo "🧪 Ejecutando tests de optimizaciones..."
	go test ./tests -v -run "TestAlumnoServiceOptimized|TestMonitoringAndConfiguration|TestPerformanceMonitoring|TestErrorHandling|TestConfigurationValidation"

# Ejecutar tests de rendimiento (benchmarks)
test-benchmark:
	@echo "📊 Ejecutando benchmarks..."
	go test ./tests -v -bench=. -benchmem

# Ejecutar todos los tests
test-all: test test-optimized test-benchmark

# Ejecutar tests con coverage
test-coverage:
	@echo "🧪 Ejecutando tests con coverage..."
	go test ./tests -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Reporte de coverage generado: coverage.html"

# Ejecutar tests con coverage de optimizaciones
test-coverage-optimized:
	@echo "🧪 Ejecutando tests de optimizaciones con coverage..."
	go test ./tests -coverprofile=coverage_optimized.out -covermode=atomic -run "TestAlumnoServiceOptimized|TestMonitoringAndConfiguration|TestPerformanceMonitoring|TestErrorHandling|TestConfigurationValidation"
	go tool cover -html=coverage_optimized.out -o coverage_optimized.html
	@echo "📊 Reporte de coverage de optimizaciones generado: coverage_optimized.html"

# Limpiar archivos generados
clean:
	@echo "🧹 Limpiando archivos..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html
	rm -f coverage_optimized.out
	rm -f coverage_optimized.html
	@echo "✅ Limpieza completada"

# Levantar servicios Docker
docker-up:
	@echo "🐳 Levantando servicios Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "✅ Servicios Docker iniciados"

# Detener servicios Docker
docker-down:
	@echo "🐳 Deteniendo servicios Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "✅ Servicios Docker detenidos"

# Reiniciar servicios Docker
docker-restart: docker-down docker-up

# Ver logs de Docker
docker-logs:
	@echo "📋 Mostrando logs de Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Verificar estado de servicios
docker-status:
	@echo "🔍 Estado de servicios Docker:"
	docker-compose -f $(DOCKER_COMPOSE_FILE) ps

# Ejecutar todo el pipeline: Docker + Build + Run
all: docker-up
	@echo "⏳ Esperando que PostgreSQL esté listo..."
	@sleep 5
	@make build
	@make run-only

# Ejecutar solo la carga masiva (asumiendo que Docker ya está corriendo)
masivo: build
	@echo "🚀 Ejecutando carga masiva de 2.5M registros..."
	./$(BINARY_NAME)

# Monitorear sistema y progreso
monitor:
	@echo "📊 Iniciando monitoreo..."
	@./scripts/monitor.sh

# Monitorear continuamente (cada 10 segundos)
monitor-continuous:
	@echo "📊 Monitoreo continuo (Ctrl+C para detener)..."
	@while true; do \
		clear; \
		./scripts/monitor.sh; \
		sleep 10; \
	done

# Pipeline completo con tests y monitoreo
pipeline: docker-up
	@echo "⏳ Esperando que PostgreSQL esté listo..."
	@sleep 5
	@echo "🧪 Ejecutando tests antes de la carga..."
	@make test-optimized
	@echo "📊 Estado inicial del sistema:"
	@make monitor
	@echo ""
	@echo "🚀 Iniciando carga masiva..."
	@make masivo &
	@echo "📊 Monitoreando progreso..."
	@make monitor-continuous

# Pipeline completo con todos los tests
pipeline-full: docker-up
	@echo "⏳ Esperando que PostgreSQL esté listo..."
	@sleep 5
	@echo "🧪 Ejecutando todos los tests..."
	@make test-all
	@echo "📊 Estado inicial del sistema:"
	@make monitor
	@echo ""
	@echo "🚀 Iniciando carga masiva..."
	@make masivo &
	@echo "📊 Monitoreando progreso..."
	@make monitor-continuous

# Instalar dependencias
deps:
	@echo "📦 Instalando dependencias..."
	go mod download
	go mod tidy
	@echo "✅ Dependencias instaladas"

# Verificar formato del código
fmt:
	@echo "🎨 Formateando código..."
	go fmt ./...
	@echo "✅ Formato aplicado"

# Verificar linting
lint:
	@echo "🔍 Verificando linting..."
	golangci-lint run
	@echo "✅ Linting completado"

# Generar documentación
docs:
	@echo "📚 Generando documentación..."
	godoc -http=:6060 &
	@echo "📖 Documentación disponible en: http://localhost:6060"

# Pipeline completo con monitoreo
pipeline: docker-up
	@echo "⏳ Esperando que PostgreSQL esté listo..."
	@sleep 5
	@echo "📊 Estado inicial del sistema:"
	@make monitor
	@echo ""
	@echo "🚀 Iniciando carga masiva..."
	@make masivo &
	@echo "📊 Monitoreando progreso..."
	@make monitor-continuous

# Ayuda
help:
	@echo "📋 Comandos disponibles:"
	@echo ""
	@echo "🔨 Compilación:"
	@echo "  build          - Compilar el programa"
	@echo "  clean          - Limpiar archivos generados"
	@echo ""
	@echo "🚀 Ejecución:"
	@echo "  run            - Compilar y ejecutar"
	@echo "  run-only       - Ejecutar sin recompilar"
	@echo "  run-test       - Ejecutar con datos de prueba"
	@echo "  masivo         - Ejecutar carga masiva de 2.5M registros"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-up      - Levantar servicios Docker"
	@echo "  docker-down    - Detener servicios Docker"
	@echo "  docker-restart - Reiniciar servicios Docker"
	@echo "  docker-logs    - Ver logs de Docker"
	@echo "  docker-status  - Estado de servicios"
	@echo ""
	@echo "📊 Monitoreo:"
	@echo "  monitor        - Monitorear sistema una vez"
	@echo "  monitor-continuous - Monitorear continuamente"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test           - Ejecutar tests básicos"
	@echo "  test-optimized - Ejecutar tests de optimizaciones"
	@echo "  test-benchmark - Ejecutar benchmarks"
	@echo "  test-all       - Ejecutar todos los tests"
	@echo "  test-coverage  - Tests con reporte de coverage"
	@echo "  test-coverage-optimized - Coverage de optimizaciones"
	@echo ""
	@echo "🛠️  Desarrollo:"
	@echo "  deps           - Instalar dependencias"
	@echo "  fmt            - Formatear código"
	@echo "  lint           - Verificar linting"
	@echo "  docs           - Generar documentación"
	@echo ""
	@echo "🎯 Pipeline completo:"
	@echo "  all            - Docker + Build + Run"
	@echo "  pipeline       - Docker + Build + Run + Monitoreo"
	@echo "  pipeline-full  - Docker + Tests + Build + Run + Monitoreo"
	@echo ""
	@echo "📋 Ejemplo de uso para 2.5M registros:"
	@echo "  make docker-up"
	@echo "  make test-optimized"
	@echo "  make masivo"
	@echo ""
	@echo "📊 Monitoreo durante la carga:"
	@echo "  make monitor-continuous" 