CREATE TYPE gender AS ENUM('male', 'female');

CREATE TYPE role AS ENUM('patient', 'doctor', 'admin');

CREATE TABLE users (
    id UUID default gen_random_uuid() PRIMARY KEY,
    email varchar(225),
    password_hash varchar(225),
    first_name varchar(225),
    last_name varchar(225),
    date_of_birth varchar(11),
    gender gender,
    role role DEFAULT 'patient',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0,
    UNIQUE(email, deleted_at)
);