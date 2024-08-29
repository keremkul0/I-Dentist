#!/bin/bash

# Remove any existing container with the same name
docker rm -f postgres-clinicRepository

# Start the PostgreSQL container
docker run --name postgres-clinicRepository -e POSTGRES_USER=clinicuser -e POSTGRES_PASSWORD=clinicpassword -p 5432:5432 -d postgres:latest

echo "Postgresql starting..."
sleep 4

# Check if the container is running
if ! docker ps | grep -q postgres-clinicRepository; then
  echo "Failed to start postgres-clinic container."
  exit 1
fi

# Create the database
docker exec -it postgres-clinicRepository psql -U clinicuser -d postgres -c "CREATE DATABASE clinicdb"
sleep 3
echo "Database clinicdb created"

# Create the tables
docker exec -it postgres-clinicRepository psql -U clinicuser -d clinicdb -c "
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    identity_number VARCHAR(20) NOT NULL UNIQUE,  
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    clinic_id BIGINT REFERENCES clinics(id)
);

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL REFERENCES users(id),
    role_id BIGINT NOT NULL REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS patients (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    identity_number VARCHAR(20) NOT NULL UNIQUE,
    birth_date DATE,
    phone_number VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    medical_history TEXT
);


CREATE TABLE IF NOT EXISTS appointments (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id),
    clinic_id BIGINT NOT NULL REFERENCES clinics(id),
    doctor_id BIGINT NOT NULL REFERENCES users(id),
    assistant_id BIGINT REFERENCES users(id),
    date_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(255),
    notes TEXT
);

CREATE TABLE IF NOT EXISTS treatments (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    clinic_id BIGINT NOT NULL REFERENCES clinics(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS appointment_treatments (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    appointment_id BIGINT NOT NULL REFERENCES appointments(id),
    treatment_id BIGINT NOT NULL REFERENCES treatments(id),
    performed_by BIGINT REFERENCES users(id),
    tooth VARCHAR(2)
);

CREATE TABLE IF NOT EXISTS clinics (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    group_id BIGINT REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_info TEXT
);
"
sleep 3
echo "Tables created"
