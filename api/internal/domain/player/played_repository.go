package player

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lardira/playtrack/internal/pkg/types"
)

const (
	TablePlayer     = "player"
	TablePlayedGame = "played_game"
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

	sqlBuild := sq.Select(playedGameColumns).
		PlaceholderFormat(sq.Dollar).
		From(TablePlayedGame).
		Where(sq.Eq{"player_id": playerID})

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, query, args...)
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
	sqlBuild := sq.Select(playedGameColumns).
		PlaceholderFormat(sq.Dollar).
		From(TablePlayedGame).
		Where(sq.Eq{"player_id": playerID, "id": id})

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.pool.QueryRow(ctx, query, args...)
	p, err := playedGameFromRow(row)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PGPlayedRepository) FindLastNotReroll(ctx context.Context, playerID string) (*PlayedGame, error) {
	sqlBuild := sq.Select(playedGameColumns).
		PlaceholderFormat(sq.Dollar).
		From(TablePlayedGame).
		Where(sq.Eq{"player_id": playerID}, sq.NotEq{"status": PlayedGameStatusRerolled}).
		OrderBy("started_at DESC").
		Limit(1).
		Offset(1)

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.pool.QueryRow(ctx, query, args...)
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

	sqlBuild := sq.Insert(TablePlayedGame).
		PlaceholderFormat(sq.Dollar).
		Columns("player_id", "game_id", "status", "points").
		Values(game.PlayerID, game.GameID, PlayedGameStatusAdded, game.Points).
		Suffix("RETURNING id")

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return id, err
	}
	row := r.pool.QueryRow(ctx, query, args...)

	if err := row.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r *PGPlayedRepository) Update(ctx context.Context, game *PlayedGameUpdate) (int, error) {
	var id int

	updBuild := sq.Update(TablePlayedGame).PlaceholderFormat(sq.Dollar)

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
		updBuild = updBuild.Set("play_time", game.PlayTime.Duration)
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
	var ptime *time.Duration
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
		&ptime,
	)
	if err != nil {
		return nil, err
	}
	if ptime != nil {
		ds := types.NewDurationString(*ptime)
		p.PlayTime = &ds
	}
	return &p, nil
}
