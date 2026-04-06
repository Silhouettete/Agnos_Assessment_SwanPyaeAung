CREATE TABLE
    IF NOT EXISTS hospitals (
        hospital_id SERIAL PRIMARY KEY,
        hospital_name TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS staff (
        staff_id SERIAL PRIMARY KEY,
        hospital_id INTEGER NOT NULL REFERENCES hospitals (hospital_id),
        first_name_th TEXT NOT NULL,
        middle_name_th TEXT NOT NULL,
        last_name_th TEXT NOT NULL,
        first_name_en TEXT NOT NULL,
        middle_name_en TEXT NOT NULL,
        last_name_en TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        role TEXT NOT NULL DEFAULT 'staff',
        created_at TIMESTAMPTZ DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS patients (
        patient_hn TEXT PRIMARY KEY,
        national_id TEXT UNIQUE,
        passport_id TEXT UNIQUE,
        first_name_th TEXT,
        middle_name_th TEXT,
        last_name_th TEXT,
        first_name_en TEXT,
        middle_name_en TEXT,
        last_name_en TEXT,
        date_of_birth DATE,
        phone_number TEXT,
        email TEXT,
        gender TEXT,
        hospital_id INTEGER REFERENCES hospitals (hospital_id)
    );

CREATE INDEX IF NOT EXISTS idx_patients_national_id ON patients (national_id);

CREATE INDEX IF NOT EXISTS idx_patients_passport_id ON patients (passport_id);

CREATE INDEX IF NOT EXISTS idx_patients_hospital_id ON patients (hospital_id);

INSERT INTO
    hospitals (hospital_name)
VALUES
    ('Bangkok General Hospital');

INSERT INTO
    patients (
        patient_hn,
        national_id,
        passport_id,
        first_name_th,
        middle_name_th,
        last_name_th,
        first_name_en,
        middle_name_en,
        last_name_en,
        date_of_birth,
        phone_number,
        email,
        gender,
        hospital_id
    )
VALUES
    (
        'HN001',
        '1234567890123',
        'PP123456',
        'สมชาย',
        NULL,
        'ใจดี',
        'Somchai',
        NULL,
        'Jaidee',
        '1990-01-15',
        '0812345678',
        'somchai@email.com',
        'Male',
        1
    );

    {"national_id":"1234567890123","passport_id":"PP123456","first_name_th":"สมชาย","middle_name_th":null,"last_name_th":"ใจดี","first_name_en":"Somchai","middle_name_en":null,"last_name_en":"Jaidee","date_of_birth":"1990-01-15","patient_hn":"HN001","phone_number":"0812345678","email":"somchai@email.com","gender":"Male","hospital_id":1}