-- +goose Up
-- +goose StatementBegin
CREATE TYPE played_game_status AS ENUM ('added', 'in_progress', 'completed', 'dropped', 'rerolled');

CREATE TABLE played_game(
    id SERIAL PRIMARY KEY,
    player_id UUID REFERENCES player(id),
    game_id INT REFERENCES game(id),
    points INT NOT NULL DEFAULT 0,
    comment TEXT NULL,
    rating INT NULL,
    status played_game_status NOT NULL DEFAULT 'added',
    play_time INTERVAL NULL,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE played_game;
DROP TYPE played_game_status;
-- +goose StatementEnd
