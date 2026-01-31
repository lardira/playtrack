package player

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	playerColumns     string = "id, username, img, email, password, created_at"
	playedGameColumns string = `id, player_id, game_id, points, comment, 
	rating, status, started_at, completed_at, play_time`
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

	query := fmt.Sprintf(`SELECT %s FROM player`, playerColumns)
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p, err := playerFromRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, nil
}

func (r *PGRepository) FindOne(ctx context.Context, id string) (*Player, error) {
	query := fmt.Sprintf(`SELECT %s FROM player WHERE id=$1`, playerColumns)
	row := r.pool.QueryRow(ctx, query, id)
	p, err := playerFromRow(row)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PGRepository) FindOneByUsername(ctx context.Context, username string) (*Player, error) {
	query := fmt.Sprintf(`SELECT %s FROM player WHERE username=$1`, playerColumns)
	row := r.pool.QueryRow(ctx, query, username)
	p, err := playerFromRow(row)
	if err != nil {
		return nil, err
	}
	return p, nil
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

func (r *PGRepository) Update(ctx context.Context, player *PlayerUpdate) (string, error) {
	var id string
	updBuild := sq.Update("player").PlaceholderFormat(sq.Dollar)

	if player.Email != nil {
		updBuild = updBuild.Set("email", *player.Email)
	}
	if player.Img != nil {
		updBuild = updBuild.Set("img", *player.Img)
	}
	if player.Username != nil {
		updBuild = updBuild.Set("username", *player.Username)
	}
	if player.Password != nil {
		updBuild = updBuild.Set("password", *player.Password)
	}

	query, args, err := updBuild.Where(sq.Eq{"id": player.ID}).Suffix("RETURNING id").ToSql()
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

func playerFromRow(row pgx.Row) (*Player, error) {
	var p Player
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
