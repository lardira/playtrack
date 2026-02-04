package game

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TableGame = "game"
)

const (
	gameColumns string = "id, points, hours_to_beat, title, url, created_at"
)

var (
	ErrFoundByTitle = errors.New("title is not unique")
)

type PGRepository struct {
	pool *pgxpool.Pool
}

func NewPGRepository(pool *pgxpool.Pool) *PGRepository {
	return &PGRepository{
		pool: pool,
	}
}

func (r *PGRepository) FindAll(ctx context.Context) ([]Game, error) {
	out := make([]Game, 0)

	sqlBuild := sq.Select(gameColumns).
		PlaceholderFormat(sq.Dollar).
		From(TableGame)

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
		g, err := gameFromRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *g)
	}
	return out, nil
}

func (r *PGRepository) FindOne(ctx context.Context, id int) (*Game, error) {
	sqlBuild := sq.Select(gameColumns).
		PlaceholderFormat(sq.Dollar).
		From(TableGame).
		Where(sq.Eq{"id": id})

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	g, err := gameFromRow(row)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *PGRepository) FindByTitle(ctx context.Context, title string) (*Game, error) {
	query := `SELECT 
				id, points, hours_to_beat, title, url, created_at  
			FROM game 
			WHERE title=$1`

	sqlBuild := sq.Select(gameColumns).
		PlaceholderFormat(sq.Dollar).
		From(TableGame).
		Where(sq.Eq{"title": title})

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	g, err := gameFromRow(row)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *PGRepository) Insert(ctx context.Context, game *Game) (int, error) {
	var id int

	sqlBuild := sq.Insert(TableGame).
		PlaceholderFormat(sq.Dollar).
		Columns("points", "hours_to_beat", "title", "url").
		Values(game.Points, game.HoursToBeat, game.Title, game.URL).
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

func gameFromRow(row pgx.Row) (*Game, error) {
	var g Game
	err := row.Scan(
		&g.ID,
		&g.Points,
		&g.HoursToBeat,
		&g.Title,
		&g.URL,
		&g.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
