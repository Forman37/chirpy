-- Get into using psql "postgres://jamesforman:@localhost:5432/chirpy"
-- Example usage : goose postgres "postgres://jamesforman:@localhost:5432/chirpy" up
-- List tables \dt inside db

-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
