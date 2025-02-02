-- +goose Up
CREATE TABLE songs2 (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date VARCHAR(255),
    text TEXT,
    lyrics VARCHAR(255),
    link VARCHAR(255)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE songs;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
