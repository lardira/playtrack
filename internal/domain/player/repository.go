package player

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRepository struct {
	pool *pgxpool.Pool
}

func NewPGRepository(pool *pgxpool.Pool) *PGRepository {
	return &PGRepository{
		pool: pool,
	}
}

func (r *PGRepository) FindAll(ctx context.Context) ([]Player, error) {
	out := make([]Player, 0)

	query := `SELECT 
				id, username, img, email, password, created_at  
			FROM player`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p Player
		err := rows.Scan(
			&p.ID,
			&p.Username,
			&p.Img,
			&p.Email,
			&p.Password,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *PGRepository) FindOne(ctx context.Context, id string) (*Player, error) {
	var p Player
	query := `SELECT 
				id, username, img, email, password, created_at  
			FROM player
			WHERE id=$1`

	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&p.ID,
		&p.Username,
		&p.Img,
		&p.Email,
		&p.Password,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PGRepository) Insert(ctx context.Context, player *Player) (string, error) {
	var id string

	query := `INSERT INTO player (id, username, img, email, password) 
			VALUES (@id, @username, @img, @email, @password)
			RETURNING id`

	args := pgx.NamedArgs{
		"id":       uuid.NewString(),
		"username": player.Username,
		"img":      player.Img,
		"email":    player.Email,
		"password": player.Password,
	}

	row := r.pool.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
