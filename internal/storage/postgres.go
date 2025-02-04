package storage

import (
	"context"
	"time"

	"github.com/erknas/song-library/internal/config"
	"github.com/erknas/song-library/internal/lib"
	"github.com/erknas/song-library/internal/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ctxTimeout time.Duration = time.Second * 5

type PostgresPool struct {
	pool *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*PostgresPool, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	poolCfg, err := lib.PoolConfig(cfg)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgresPool{pool: pool}, nil
}

func (p *PostgresPool) Songs(ctx context.Context, pag types.Pagination) ([]*types.Song, error) {
	query := `SELECT id, song, group_name, release_date, text, link 
		  	  FROM songs
			  ORDER BY id 
			  LIMIT @size 
			  OFFSET @page
			 `
	args := pgx.NamedArgs{
		"size": pag.Size,
		"page": pag.Page,
	}

	rows, err := p.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*types.Song

	for rows.Next() {
		song := new(types.Song)
		if err := rows.Scan(&song.ID, &song.Song, &song.Group, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (p *PostgresPool) SongsByFilters(ctx context.Context, query string, args []any) ([]*types.Song, error) {
	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*types.Song

	for rows.Next() {
		song := new(types.Song)
		if err := rows.Scan(&song.ID, &song.Song, &song.Group, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (p *PostgresPool) SongText(ctx context.Context, id int) (string, error) {
	query := `SELECT text 
			  FROM songs 
			  WHERE id=@id
			 `

	args := pgx.NamedArgs{
		"id": id,
	}

	row := p.pool.QueryRow(ctx, query, args)

	var text string

	if err := row.Scan(&text); err != nil {
		return text, err
	}

	return text, nil
}

func (p *PostgresPool) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs 
			  WHERE id=@id
			 `

	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := p.pool.Exec(ctx, query, args)

	return err
}

func (p *PostgresPool) UpdateSong(ctx context.Context, id int, song *types.Song) error {
	query := `UPDATE songs 
			  SET
			  song=@song,
			  group_name=@group_name, 
			  release_date=@release_date, 
			  text=@text, 
			  link=@link
			  WHERE id=@id
			 `

	args := pgx.NamedArgs{
		"song":         song.Song,
		"group_name":   song.Group,
		"release_date": song.ReleaseDate,
		"text":         song.Text,
		"link":         song.Link,
		"id":           id,
	}

	_, err := p.pool.Exec(ctx, query, args)

	return err
}

func (p *PostgresPool) AddSong(ctx context.Context, song *types.Song) error {
	query := `INSERT INTO songs(song, group_name, release_date, text, link)
			  VALUES(@song, @group_name, @release_date, @text, @link)
			 `
	args := pgx.NamedArgs{
		"song":         song.Song,
		"group_name":   song.Group,
		"release_date": song.ReleaseDate,
		"text":         song.Text,
		"link":         song.Link,
	}

	_, err := p.pool.Exec(ctx, query, args)

	return err
}

func (p *PostgresPool) Close() {
	p.pool.Close()
}
