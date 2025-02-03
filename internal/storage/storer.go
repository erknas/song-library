package storage

import (
	"context"

	"github.com/erknas/song-library/internal/types"
)

type Storer interface {
	Songs(context.Context, types.Pagination) ([]*types.Song, error)
	SongsByFilters(context.Context, types.Pagination, types.Filter) ([]*types.Song, error)
	SongText(context.Context, types.Pagination, int) (string, error)
	DeleteSong(context.Context, int) error
	UpdateSong(context.Context, int, *types.Song) error
	AddSong(context.Context, *types.Song) error
}
