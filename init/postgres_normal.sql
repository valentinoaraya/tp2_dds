-- Script para restaurar configuración normal de PostgreSQL después de carga masiva
-- Este script debe ejecutarse después de la carga masiva

-- Restaurar configuraciones de rendimiento
SET synchronous_commit = on;            -- Habilitar commits síncronos
SET wal_buffers = '4MB';                -- Buffer WAL normal
SET checkpoint_segments = 3;             -- Segmentos de checkpoint normales
SET checkpoint_completion_target = 0.5;  -- Checkpoints normales
SET wal_writer_delay = '200ms';          -- Escritor WAL normal
SET commit_delay = 0;                    -- Sin delay en commits
SET commit_siblings = 5;                 -- Siblings normales

-- Restaurar configuraciones de memoria
SET shared_buffers = '128MB';            -- Buffer compartido normal
SET effective_cache_size = '256MB';      -- Cache efectivo normal
SET work_mem = '4MB';                    -- Memoria de trabajo normal
SET maintenance_work_mem = '64MB';       -- Memoria de mantenimiento normal

-- Restaurar configuraciones de conexiones
SET max_connections = 100;               -- Conexiones normales
SET max_prepared_transactions = 0;       -- Sin transacciones preparadas

-- Restaurar configuraciones de logging
SET log_statement = 'none';              -- No loggear statements
SET log_min_duration_statement = -1;     -- No loggear por duración

-- Restaurar configuraciones de autovacuum
SET autovacuum = on;                     -- Habilitar autovacuum
SET track_counts = on;                   -- Trackear estadísticas

-- Restaurar configuraciones de planificación
SET random_page_cost = 4.0;              -- Costo de página aleatoria normal
SET effective_io_concurrency = 1;        -- I/O concurrente normal

-- Restaurar configuraciones de locks
SET deadlock_timeout = '1s';             -- Timeout de deadlock normal
SET lock_timeout = '0';                  -- Sin timeout de locks

-- Restaurar configuraciones de archivos temporales
SET temp_file_limit = -1;                -- Sin límite de archivos temporales

-- Restaurar configuraciones de estadísticas
SET track_activities = on;               -- Trackear actividades
SET track_io_timing = off;               -- No trackear I/O timing

-- Restaurar configuraciones de WAL
SET wal_level = replica;                 -- Nivel WAL normal
SET max_wal_senders = 10;                -- Senders WAL normales

-- Restaurar configuraciones de archivos
SET max_files_per_process = 1000;        -- Archivos por proceso normales

-- Restaurar configuraciones de timeout
SET statement_timeout = '0';             -- Sin timeout de statements
SET idle_in_transaction_session_timeout = '0'; -- Sin timeout de sesiones idle

-- Restaurar configuraciones de buffer
SET shared_preload_libraries = '';       -- Sin librerías pre-cargadas

-- Restaurar configuraciones de planificación de consultas
SET enable_seqscan = on;                 -- Habilitar sequential scan
SET enable_indexscan = on;               -- Habilitar index scan
SET enable_bitmapscan = on;              -- Habilitar bitmap scan
SET enable_hashjoin = on;                -- Habilitar hash join
SET enable_mergejoin = on;               -- Habilitar merge join
SET enable_nestloop = on;                -- Habilitar nested loop

-- Restaurar configuraciones de costos
SET cpu_tuple_cost = 0.01;               -- Costo de CPU por tupla normal
SET cpu_index_tuple_cost = 0.005;        -- Costo de CPU por índice normal
SET cpu_operator_cost = 0.0025;          -- Costo de CPU por operador normal

-- Restaurar configuraciones de paralelismo
SET max_parallel_workers_per_gather = 2; -- Workers paralelos normales
SET max_parallel_workers = 8;            -- Workers paralelos totales normales
SET parallel_tuple_cost = 0.1;           -- Costo de tupla paralela normal
SET parallel_setup_cost = 1000.0;        -- Costo de setup paralelo normal

-- Restaurar configuraciones de JIT
SET jit = on;                            -- Habilitar JIT

-- Restaurar configuraciones de logging de errores
SET log_min_messages = 'warning';        -- Mensajes de warning en log
SET log_line_prefix = '%t ';             -- Prefijo de log normal

-- Restaurar configuraciones de estadísticas de planificación
SET track_functions = 'pl';              -- Trackear funciones PL

-- Restaurar configuraciones de archivos de datos
SET data_directory = '';                 -- Directorio de datos por defecto

-- Restaurar configuraciones de archivos de configuración
SET hba_file = '';                       -- Archivo HBA por defecto
SET ident_file = '';                     -- Archivo ident por defecto

-- Restaurar configuraciones de SSL
SET ssl = off;                           -- SSL deshabilitado

-- Restaurar configuraciones de autenticación
SET password_encryption = 'md5';         -- Encriptación de contraseñas

-- Restaurar configuraciones de timezone
SET timezone = 'UTC';                    -- Timezone UTC

-- Restaurar configuraciones de locale
SET lc_messages = 'C';                   -- Locale de mensajes
SET lc_monetary = 'C';                   -- Locale monetario
SET lc_numeric = 'C';                    -- Locale numérico
SET lc_time = 'C';                       -- Locale de tiempo

-- Restaurar configuraciones de encoding
SET client_encoding = 'UTF8';            -- Encoding del cliente
SET server_encoding = 'UTF8';            -- Encoding del servidor

-- Restaurar configuraciones de collation
SET lc_collate = 'C';                    -- Collation
SET lc_ctype = 'C';                      -- Ctype

-- Ejecutar VACUUM para limpiar y actualizar estadísticas
VACUUM ANALYZE;

-- Mostrar configuración restaurada
SELECT 'Configuración normal de PostgreSQL restaurada' as status; 