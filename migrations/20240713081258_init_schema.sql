-- +goose Up
CREATE TABLE IF NOT EXISTS tbl_users (
    id SERIAL PRIMARY KEY ,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    age INTEGER NOT NULL,
    password TEXT NOT NULL,

    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS tbl_attendances (
    id SERIAL PRIMARY KEY ,
    user_id INTEGER NOT NULL REFERENCES tbl_users(id),
    date DATE NOT NULL,
    time_in TIME NULL,
    time_out TIME NULL,
    status VARCHAR(10) NOT NULL CHECK(status IN ('present', 'absent', 'late', 'sick')),
    notes TEXT,
    
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    deleted_at timestamp
);




-- +goose Down
DROP TABLE IF EXISTS tbl_users cascade;
DROP TABLE IF EXISTS tbl_attendances cascade;
