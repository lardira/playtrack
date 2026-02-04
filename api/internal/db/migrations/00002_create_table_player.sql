-- +goose Up
-- +goose StatementBegin
CREATE TABLE player(
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    img TEXT NULL,
    email TEXT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE player;
-- +goose StatementEnd
