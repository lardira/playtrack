package player

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPlayedGameNotFound = errors.New("played game is not found")
)

type PGPlayedRepository struct {
	pool *pgxpool.Pool
}

func NewPGPlayedRepository(pool *pgxpool.Pool) *PGPlayedRepository {
	return &PGPlayedRepository{
		pool: pool,
	}
}

func (r *PGPlayedRepository) FindAll(ctx context.Context, playerID string) ([]PlayedGame, error) {
	out := make([]PlayedGame, 0)

	query := fmt.Sprintf(`SELECT %s FROM played_game WHERE player_id=$1`, playedGameColumns)
	rows, err := r.pool.Query(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p, err := playedGameFromRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, nil
}

func (r *PGPlayedRepository) FindOne(ctx context.Context, playerID string, id int) (*PlayedGame, error) {
	query := fmt.Sprintf(`SELECT %s FROM played_game WHERE player_id=$1 AND id=$2`, playedGameColumns)
	row := r.pool.QueryRow(ctx, query, playerID, id)
	p, err := playedGameFromRow(row)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PGPlayedRepository) FindLast(ctx context.Context, playerID string) (*PlayedGame, error) {
	query := fmt.Sprintf(`SELECT %s FROM played_game 
		WHERE player_id=$1 
		ORDER BY started_at DESC LIMIT 1 OFFSET 1`, playedGameColumns)
	row := r.pool.QueryRow(ctx, query, playerID)
	p, err := playedGameFromRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlayedGameNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *PGPlayedRepository) Insert(ctx context.Context, game *PlayedGame) (int, error) {
	var id int

	query := `INSERT INTO played_game (player_id, game_id, status, points) 
			VALUES (@player_id, @game_id, @status, @points)
			RETURNING id`

	args := pgx.NamedArgs{
		"player_id": game.PlayerID,
		"game_id":   game.GameID,
		"status":    PlayedGameStatusAdded,
		"points":    game.Points,
	}

	row := r.pool.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *PGPlayedRepository) Update(ctx context.Context, game *PlayedGameUpdate) (int, error) {
	var id int
	updBuild := sq.Update("played_game").PlaceholderFormat(sq.Dollar)

	if game.Points != nil {
		updBuild = updBuild.Set("points", *game.Points)
	}
	if game.Comment != nil {
		updBuild = updBuild.Set("comment", *game.Comment)
	}
	if game.Rating != nil {
		updBuild = updBuild.Set("rating", *game.Rating)
	}
	if game.Status != nil {
		updBuild = updBuild.Set("status", *game.Status)
	}
	if game.CompletedAt != nil {
		updBuild = updBuild.Set("completed_at", *game.CompletedAt)
	}
	if game.PlayTime != nil {
		updBuild = updBuild.Set("play_time", *game.PlayTime)
	}

	query, args, err := updBuild.Where(sq.Eq{"id": game.ID}).Suffix("RETURNING id").ToSql()
	if err != nil {
		return id, err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func playedGameFromRow(row pgx.Row) (*PlayedGame, error) {
	var p PlayedGame
	err := row.Scan(
		&p.ID,
		&p.PlayerID,
		&p.GameID,
		&p.Points,
		&p.Comment,
		&p.Rating,
		&p.Status,
		&p.StartedAt,
		&p.CompletedAt,
		&p.PlayTime,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
