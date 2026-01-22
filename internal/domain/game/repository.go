package game

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
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

	query := `SELECT 
				id, points, hours_to_beat, title, url, created_at  
			FROM game`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var g Game
		err := rows.Scan(
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
		out = append(out, g)
	}
	return out, nil
}

func (r *PGRepository) FindOne(ctx context.Context, id int) (*Game, error) {
	var game Game

	query := `SELECT 
				id, points, hours_to_beat, title, url, created_at  
			FROM game 
			WHERE id=$1`

	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&game.ID,
		&game.Points,
		&game.HoursToBeat,
		&game.Title,
		&game.URL,
		&game.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *PGRepository) FindByTitle(ctx context.Context, title string) (*Game, error) {
	var game Game

	query := `SELECT 
				id, points, hours_to_beat, title, url, created_at  
			FROM game 
			WHERE title=$1`

	row := r.pool.QueryRow(ctx, query, title)
	err := row.Scan(
		&game.ID,
		&game.Points,
		&game.HoursToBeat,
		&game.Title,
		&game.URL,
		&game.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *PGRepository) Insert(ctx context.Context, game *Game) (int, error) {
	if found, _ := r.FindByTitle(ctx, game.Title); found != nil {
		return 0, ErrFoundByTitle
	}

	var id int

	query := `INSERT INTO 
				game (points, hours_to_beat, title, url) 
			VALUES ($1, $2, $3, $4)
			RETURNING id`

	row := r.pool.QueryRow(ctx, query, game.Points, game.HoursToBeat, game.Title, game.URL)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
