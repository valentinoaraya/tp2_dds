-- Script de optimización de PostgreSQL para carga masiva de 2.5M registros
-- Este script debe ejecutarse antes de la carga masiva

-- Configuraciones de rendimiento para carga masiva
SET synchronous_commit = off;           -- Reducir latencia de escritura
SET wal_buffers = '16MB';               -- Buffer WAL más grande
SET checkpoint_segments = 32;           -- Más segmentos de checkpoint
SET checkpoint_completion_target = 0.9; -- Checkpoints más suaves
SET wal_writer_delay = '200ms';         -- Escritor WAL más frecuente
SET commit_delay = 1000;                -- Agrupar commits
SET commit_siblings = 5;                -- Mínimo de transacciones concurrentes

-- Configuraciones de memoria
SET shared_buffers = '256MB';           -- Buffer compartido
SET effective_cache_size = '1GB';       -- Cache efectivo
SET work_mem = '4MB';                   -- Memoria de trabajo por operación
SET maintenance_work_mem = '64MB';      -- Memoria para mantenimiento

-- Configuraciones de conexiones
SET max_connections = 100;              -- Más conexiones
SET max_prepared_transactions = 100;    -- Transacciones preparadas

-- Configuraciones de logging (reducir durante carga)
SET log_statement = 'none';             -- No loggear statements
SET log_min_duration_statement = 1000;  -- Solo loggear queries lentas (>1s)

-- Configuraciones de autovacuum (deshabilitar durante carga)
SET autovacuum = off;                   -- Deshabilitar autovacuum
SET track_counts = off;                 -- No trackear estadísticas

-- Configuraciones de planificación
SET random_page_cost = 1.1;             -- Costo de página aleatoria (SSD)
SET effective_io_concurrency = 200;     -- I/O concurrente

-- Configuraciones de locks
SET deadlock_timeout = '1s';            -- Timeout de deadlock corto
SET lock_timeout = '30s';               -- Timeout de locks

-- Configuraciones de archivos temporales
SET temp_file_limit = '1GB';            -- Límite de archivos temporales

-- Configuraciones de estadísticas
SET track_activities = off;             -- No trackear actividades
SET track_io_timing = off;              -- No trackear I/O timing

-- Configuraciones de WAL
SET wal_level = minimal;                -- Nivel mínimo de WAL
SET max_wal_senders = 0;                -- Sin replicación

-- Configuraciones de archivos
SET max_files_per_process = 1000;       -- Más archivos por proceso

-- Configuraciones de timeout
SET statement_timeout = '1h';           -- Timeout de statements
SET idle_in_transaction_session_timeout = '1h'; -- Timeout de sesiones idle

-- Configuraciones de buffer
SET shared_preload_libraries = '';      -- Sin librerías pre-cargadas

-- Configuraciones de planificación de consultas
SET enable_seqscan = on;                -- Habilitar sequential scan
SET enable_indexscan = on;              -- Habilitar index scan
SET enable_bitmapscan = on;             -- Habilitar bitmap scan
SET enable_hashjoin = on;               -- Habilitar hash join
SET enable_mergejoin = on;              -- Habilitar merge join
SET enable_nestloop = on;               -- Habilitar nested loop

-- Configuraciones de costos
SET cpu_tuple_cost = 0.01;              -- Costo de CPU por tupla
SET cpu_index_tuple_cost = 0.005;       -- Costo de CPU por índice
SET cpu_operator_cost = 0.0025;         -- Costo de CPU por operador

-- Configuraciones de paralelismo
SET max_parallel_workers_per_gather = 4; -- Workers paralelos por gather
SET max_parallel_workers = 8;           -- Workers paralelos totales
SET parallel_tuple_cost = 0.1;          -- Costo de tupla paralela
SET parallel_setup_cost = 1000.0;       -- Costo de setup paralelo

-- Configuraciones de JIT (Just-In-Time compilation)
SET jit = off;                          -- Deshabilitar JIT durante carga

-- Configuraciones de archivos de configuración
SET config_file = '';                   -- Sin archivo de configuración externo

-- Configuraciones de logging de errores
SET log_min_messages = 'error';         -- Solo errores en log
SET log_line_prefix = '';               -- Sin prefijo en logs

-- Configuraciones de estadísticas de planificación
SET track_functions = 'none';           -- No trackear funciones

-- Configuraciones de archivos de datos
SET data_directory = '';                -- Directorio de datos por defecto

-- Configuraciones de archivos de configuración
SET hba_file = '';                      -- Archivo HBA por defecto
SET ident_file = '';                    -- Archivo ident por defecto

-- Configuraciones de SSL
SET ssl = off;                          -- SSL deshabilitado

-- Configuraciones de autenticación
SET password_encryption = 'md5';        -- Encriptación de contraseñas

-- Configuraciones de timezone
SET timezone = 'UTC';                   -- Timezone UTC

-- Configuraciones de locale
SET lc_messages = 'C';                  -- Locale de mensajes
SET lc_monetary = 'C';                  -- Locale monetario
SET lc_numeric = 'C';                   -- Locale numérico
SET lc_time = 'C';                      -- Locale de tiempo

-- Configuraciones de encoding
SET client_encoding = 'UTF8';           -- Encoding del cliente
SET server_encoding = 'UTF8';           -- Encoding del servidor

-- Configuraciones de collation
SET lc_collate = 'C';                   -- Collation
SET lc_ctype = 'C';                     -- Ctype

-- Configuraciones de archivos de configuración
SET config_file = '';                   -- Sin archivo de configuración externo

-- Configuraciones de logging de errores
SET log_min_messages = 'error';         -- Solo errores en log
SET log_line_prefix = '';               -- Sin prefijo en logs

-- Configuraciones de estadísticas de planificación
SET track_functions = 'none';           -- No trackear funciones

-- Configuraciones de archivos de datos
SET data_directory = '';                -- Directorio de datos por defecto

-- Configuraciones de archivos de configuración
SET hba_file = '';                      -- Archivo HBA por defecto
SET ident_file = '';                    -- Archivo ident por defecto

-- Configuraciones de SSL
SET ssl = off;                          -- SSL deshabilitado

-- Configuraciones de autenticación
SET password_encryption = 'md5';        -- Encriptación de contraseñas

-- Configuraciones de timezone
SET timezone = 'UTC';                   -- Timezone UTC

-- Configuraciones de locale
SET lc_messages = 'C';                  -- Locale de mensajes
SET lc_monetary = 'C';                  -- Locale monetario
SET lc_numeric = 'C';                   -- Locale numérico
SET lc_time = 'C';                      -- Locale de tiempo

-- Configuraciones de encoding
SET client_encoding = 'UTF8';           -- Encoding del cliente
SET server_encoding = 'UTF8';           -- Encoding del servidor

-- Configuraciones de collation
SET lc_collate = 'C';                   -- Collation
SET lc_ctype = 'C';                     -- Ctype

-- Mostrar configuración aplicada
SELECT 'Configuración de PostgreSQL optimizada para carga masiva aplicada' as status; 