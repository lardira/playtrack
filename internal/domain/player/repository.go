package player

import (
	"context"

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

	sqlBuild := sq.Select(playerColumns).
		PlaceholderFormat(sq.Dollar).
		From(TablePlayer)

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
		p, err := playerFromRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, nil
}

func (r *PGRepository) FindOne(ctx context.Context, id string) (*Player, error) {
	sqlBuild := sq.Select(playerColumns).
		PlaceholderFormat(sq.Dollar).
		From(TablePlayer).
		Where(sq.Eq{"id": id})

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	p, err := playerFromRow(row)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PGRepository) FindOneByUsername(ctx context.Context, username string) (*Player, error) {
	sqlBuild := sq.Select(playerColumns).
		PlaceholderFormat(sq.Dollar).
		From(TablePlayer).
		Where(sq.Eq{"username": username})

	query, args, err := sqlBuild.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)
	p, err := playerFromRow(row)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PGRepository) Insert(ctx context.Context, player *Player) (string, error) {
	var id string

	sqlBuild := sq.Insert(TablePlayer).
		PlaceholderFormat(sq.Dollar).
		Columns("id", "username", "img", "email", "password").
		Values(
			uuid.NewString(),
			player.Username,
			player.Img,
			player.Email,
			player.Password,
		).
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

func (r *PGRepository) Update(ctx context.Context, player *PlayerUpdate) (string, error) {
	var id string
	updBuild := sq.Update(TablePlayer).PlaceholderFormat(sq.Dollar)

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
