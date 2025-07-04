CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS file_processes (
    id UUID PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    received_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(64) NOT NULL,
    error_msg TEXT
);

CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    address VARCHAR(255),
    cnpj VARCHAR(50),
    deleted_at TIMESTAMP,
    CONSTRAINT idx_email_unique UNIQUE (email)
);

-- Índice único parcial para CNPJ preenchido
CREATE UNIQUE INDEX IF NOT EXISTS idx_cnpj_unique ON clients (cnpj) WHERE cnpj IS NOT NULL AND cnpj <> '';