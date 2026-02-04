-- +goose Up
-- +goose StatementBegin
ALTER TABLE player
    ADD COLUMN description TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE player
    DROP COLUMN description;
-- +goose StatementEnd
