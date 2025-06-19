#!/bin/bash

# Script de monitoreo para carga masiva de 2.5M registros
# Uso: ./scripts/monitor.sh

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar información del sistema
show_system_info() {
    echo -e "${BLUE}📊 INFORMACIÓN DEL SISTEMA${NC}"
    echo "=================================="
    
    # CPU
    echo -e "${GREEN}💻 CPU:${NC}"
    echo "  Cores: $(nproc)"
    echo "  Modelo: $(grep 'model name' /proc/cpuinfo | head -1 | cut -d: -f2 | xargs)"
    echo "  Uso: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)%"
    
    # Memoria
    echo -e "${GREEN}🧠 Memoria:${NC}"
    free -h | grep -E "Mem|Swap" | while read line; do
        echo "  $line"
    done
    
    # Disco
    echo -e "${GREEN}💾 Disco:${NC}"
    df -h / | tail -1 | awk '{print "  Total: " $2 " | Usado: " $3 " | Libre: " $4 " | Uso: " $5}'
    
    echo ""
}

# Función para verificar estado de Docker
check_docker_status() {
    echo -e "${BLUE}🐳 ESTADO DE DOCKER${NC}"
    echo "====================="
    
    if docker ps | grep -q postgres; then
        echo -e "${GREEN}✅ PostgreSQL está corriendo${NC}"
        
        # Información del contenedor
        CONTAINER_ID=$(docker ps | grep postgres | awk '{print $1}')
        echo "  Contenedor ID: $CONTAINER_ID"
        echo "  Puerto: $(docker port $CONTAINER_ID | grep 5432)"
        
        # Uso de recursos del contenedor
        echo -e "${YELLOW}📈 Recursos del contenedor:${NC}"
        docker stats --no-stream $CONTAINER_ID | tail -1 | awk '{print "  CPU: " $2 " | Memoria: " $3 " | Red I/O: " $4 " | Disco I/O: " $5}'
        
    else
        echo -e "${RED}❌ PostgreSQL no está corriendo${NC}"
        echo "  Ejecuta: make docker-up"
    fi
    
    echo ""
}

# Función para verificar conexión a la base de datos
check_db_connection() {
    echo -e "${BLUE}🔌 CONEXIÓN A BASE DE DATOS${NC}"
    echo "==============================="
    
    # Verificar si existe el archivo .env
    if [ ! -f .env ]; then
        echo -e "${RED}❌ Archivo .env no encontrado${NC}"
        return
    fi
    
    # Cargar variables de entorno
    source .env
    
    # Intentar conexión
    if command -v psql &> /dev/null; then
        if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $CONTAINER_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" &> /dev/null; then
            echo -e "${GREEN}✅ Conexión exitosa${NC}"
            echo "  Host: $DB_HOST"
            echo "  Puerto: $CONTAINER_PORT"
            echo "  Base de datos: $DB_NAME"
            echo "  Usuario: $DB_USER"
        else
            echo -e "${RED}❌ Error de conexión${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  psql no está instalado${NC}"
    fi
    
    echo ""
}

# Función para verificar progreso de la carga
check_load_progress() {
    echo -e "${BLUE}📈 PROGRESO DE CARGA${NC}"
    echo "====================="
    
    # Verificar si existe el archivo .env
    if [ ! -f .env ]; then
        echo -e "${RED}❌ Archivo .env no encontrado${NC}"
        return
    fi
    
    # Cargar variables de entorno
    source .env
    
    # Contar registros en la tabla
    if command -v psql &> /dev/null; then
        COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $CONTAINER_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM alumnos;" 2>/dev/null | xargs)
        
        if [ "$COUNT" != "" ]; then
            echo -e "${GREEN}✅ Registros en tabla alumnos: $COUNT${NC}"
            
            # Calcular progreso si se conoce el total esperado
            if [ "$COUNT" -gt 0 ]; then
                # Asumiendo 2.5M registros como objetivo
                TOTAL_EXPECTED=2500000
                PERCENTAGE=$((COUNT * 100 / TOTAL_EXPECTED))
                echo "  Progreso: $PERCENTAGE% ($COUNT / $TOTAL_EXPECTED)"
                
                if [ "$PERCENTAGE" -ge 100 ]; then
                    echo -e "${GREEN}🎉 ¡Carga completada!${NC}"
                fi
            fi
        else
            echo -e "${YELLOW}⚠️  No se pudo obtener el conteo${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  psql no está instalado${NC}"
    fi
    
    echo ""
}

# Función para verificar tamaño del archivo CSV
check_csv_size() {
    echo -e "${BLUE}📄 ARCHIVO CSV${NC}"
    echo "============="
    
    CSV_FILE="data/alumnos.csv"
    
    if [ -f "$CSV_FILE" ]; then
        SIZE=$(du -h "$CSV_FILE" | cut -f1)
        LINES=$(wc -l < "$CSV_FILE")
        
        echo -e "${GREEN}✅ Archivo encontrado: $CSV_FILE${NC}"
        echo "  Tamaño: $SIZE"
        echo "  Líneas: $LINES"
        
        # Calcular registros (excluyendo header)
        RECORDS=$((LINES - 1))
        echo "  Registros: $RECORDS"
        
    else
        echo -e "${RED}❌ Archivo CSV no encontrado: $CSV_FILE${NC}"
    fi
    
    echo ""
}

# Función para verificar procesos del programa
check_program_processes() {
    echo -e "${BLUE}🔄 PROCESOS DEL PROGRAMA${NC}"
    echo "========================="
    
    # Buscar procesos de Go
    GO_PROCESSES=$(ps aux | grep -E "alumnos-loader|go run" | grep -v grep | wc -l)
    
    if [ "$GO_PROCESSES" -gt 0 ]; then
        echo -e "${GREEN}✅ Procesos de Go ejecutándose: $GO_PROCESSES${NC}"
        
        # Mostrar detalles de los procesos
        ps aux | grep -E "alumnos-loader|go run" | grep -v grep | while read line; do
            PID=$(echo $line | awk '{print $2}')
            CPU=$(echo $line | awk '{print $3}')
            MEM=$(echo $line | awk '{print $4}')
            CMD=$(echo $line | awk '{for(i=11;i<=NF;i++) printf "%s ", $i; print ""}')
            echo "  PID: $PID | CPU: $CPU% | Mem: $MEM% | $CMD"
        done
    else
        echo -e "${YELLOW}⚠️  No hay procesos de Go ejecutándose${NC}"
    fi
    
    echo ""
}

# Función para mostrar logs recientes
show_recent_logs() {
    echo -e "${BLUE}📋 LOGS RECIENTES${NC}"
    echo "================"
    
    # Mostrar últimos logs de Docker
    echo -e "${YELLOW}🐳 Logs de PostgreSQL (últimas 5 líneas):${NC}"
    docker logs --tail 5 $(docker ps | grep postgres | awk '{print $1}') 2>/dev/null || echo "  No se pudieron obtener logs"
    
    echo ""
}

# Función para mostrar recomendaciones
show_recommendations() {
    echo -e "${BLUE}💡 RECOMENDACIONES${NC}"
    echo "=================="
    
    # Verificar memoria disponible
    MEM_AVAILABLE=$(free | grep Mem | awk '{print $7}')
    MEM_AVAILABLE_MB=$((MEM_AVAILABLE / 1024))
    
    if [ "$MEM_AVAILABLE_MB" -lt 1000 ]; then
        echo -e "${YELLOW}⚠️  Memoria disponible baja: ${MEM_AVAILABLE_MB}MB${NC}"
        echo "  Considera usar estrategia 'Conservador'"
    else
        echo -e "${GREEN}✅ Memoria disponible: ${MEM_AVAILABLE_MB}MB${NC}"
        echo "  Puedes usar estrategia 'Ultra-Rápido'"
    fi
    
    # Verificar CPU
    CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
    if [ "$CPU_USAGE" -gt 80 ]; then
        echo -e "${YELLOW}⚠️  Uso de CPU alto: ${CPU_USAGE}%${NC}"
        echo "  Considera reducir número de workers"
    else
        echo -e "${GREEN}✅ Uso de CPU: ${CPU_USAGE}%${NC}"
    fi
    
    echo ""
}

# Función principal
main() {
    echo -e "${BLUE}🚀 MONITOR DE CARGA MASIVA - 2.5M REGISTROS${NC}"
    echo "================================================"
    echo ""
    
    show_system_info
    check_docker_status
    check_db_connection
    check_csv_size
    check_load_progress
    check_program_processes
    show_recent_logs
    show_recommendations
    
    echo -e "${GREEN}✅ Monitoreo completado${NC}"
}

# Ejecutar función principal
main "$@" 