package service

import (
	"context"

	"github.com/erknas/song-library/internal/types"
)

type Servicer interface {
	GetSongs(context.Context, types.Pagination, types.Filter) ([]*types.Song, error)
	GetSongText(context.Context, types.Pagination, int) (string, error)
	DeleteSong(context.Context, int) error
	UpdateSong(context.Context, int, *types.Song) error
	AddSong(context.Context, *types.SongRequest) error
}
