package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/erknas/song-library/internal/lib"
	"github.com/erknas/song-library/internal/storage"
	"github.com/erknas/song-library/internal/types"
)

type Response struct {
	types.Song
	err error
}

type Service struct {
	url   string
	store storage.Storer
}

func New(url string, store storage.Storer) *Service {
	return &Service{
		url:   url,
		store: store,
	}
}

func (s *Service) GetSongs(ctx context.Context, pag types.Pagination, fil types.Filter) ([]*types.Song, error) {
	page := (pag.Page - 1) * pag.Size

	pag = types.Pagination{
		Page: page,
		Size: pag.Size,
	}

	if fil.Song == "" && fil.Group == "" && fil.Date == nil {
		return s.store.Songs(ctx, pag)
	}

	return s.store.SongsByFilters(ctx, pag, fil)
}

func (s *Service) GetSongText(ctx context.Context, pag types.Pagination, id int) (string, error) {
	return s.store.SongText(ctx, pag, id)
}

func (s *Service) DeleteSong(ctx context.Context, id int) error {
	return s.store.DeleteSong(ctx, id)
}

func (s *Service) UpdateSong(ctx context.Context, id int, req *types.Song) error {
	return s.store.UpdateSong(ctx, id, req)
}

func (s *Service) AddSong(ctx context.Context, req *types.SongRequest) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	respch := make(chan Response)

	go func() {
		song, err := s.fetchSongDetails(req)
		if err != nil {
			respch <- Response{err: err}
			return
		}
		respch <- Response{
			Song: *song,
			err:  err,
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context timeout")
	case resp := <-respch:
		if resp.err != nil {
			return resp.err
		}
		if err := s.store.AddSong(ctx, &resp.Song); err != nil {
			return err
		}
		return nil
	}
}

func (s *Service) fetchSongDetails(req *types.SongRequest) (*types.Song, error) {
	u, err := parseURL(s.url, req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	details := new(types.Details)

	if err := json.NewDecoder(resp.Body).Decode(details); err != nil {
		return nil, err
	}

	releaseDate, err := time.Parse(lib.Layout, details.ReleaseDate)
	if err != nil {
		return nil, err
	}

	song := &types.Song{
		Song:        req.Song,
		Group:       req.Group,
		ReleaseDate: releaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	return song, nil
}

func parseURL(lurl string, req *types.SongRequest) (string, error) {
	baseURL, err := url.Parse(lurl)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("group", req.Group)
	params.Add("song", req.Song)

	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}
