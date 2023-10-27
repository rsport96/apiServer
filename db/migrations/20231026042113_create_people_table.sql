-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS people (
id serial PRIMARY KEY,
name VARCHAR (50) NOT NULL,
surname VARCHAR (50) NOT NULL,
patronymic VARCHAR (50) NOT NULL,
age INT NOT NULL,
sex VARCHAR (50) NOT NULL,
country VARCHAR (50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE people;
-- +goose StatementEnd
