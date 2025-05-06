-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    birthday DATE NOT NULL,
    hash_password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE user
