-- +goose Up
-- +goose StatementBegin
CREATE TABLE game(
    id SERIAL PRIMARY KEY,
    points INT NOT NULL DEFAULT 0,
    hours_to_beat INT NOT NULL DEFAULT 0,
    title TEXT NOT NULL UNIQUE,
    url TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE game;
-- +goose StatementEnd
