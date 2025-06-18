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