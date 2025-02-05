package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/erknas/song-library/internal/errs"
	"github.com/erknas/song-library/internal/lib"
	"github.com/erknas/song-library/internal/logger"
	"github.com/erknas/song-library/internal/logger/sl"
	"github.com/erknas/song-library/internal/storage"
	"github.com/erknas/song-library/internal/types"
)

const (
	fnName             = "func"
	getSongsFn         = "GetSongs"
	getSongTextFn      = "GetSongText"
	deleteSongFn       = "DeleteSong"
	updateSongFn       = "UpdateSong"
	addSongFn          = "AddSong"
	fetchSongDetailsFn = "fetchSongDetails"
)

type Response struct {
	types.Song
	err error
}

type Service struct {
	url   string
	log   *slog.Logger
	store storage.Storer
}

func New(url string, log *slog.Logger, store storage.Storer) *Service {
	return &Service{
		url:   url,
		log:   log,
		store: store,
	}
}

func (s *Service) GetSongs(ctx context.Context, pag types.Pagination, fil types.Filter) ([]*types.Song, error) {
	log := s.log.With(slog.String(fnName, getSongsFn))

	page := (pag.Page - 1) * pag.Size

	pag = types.Pagination{
		Page: page,
		Size: pag.Size,
	}

	log.DebugContext(ctx, "pagination params", "pagination", pag)

	if fil.Song == "" && fil.Group == "" && fil.Date == nil {
		songs, err := s.store.Songs(ctx, pag)
		if err != nil {
			log.ErrorContext(ctx, "failed to get songs from w/o filters", sl.Err(err))
			return nil, err
		}

		if len(songs) == 0 {
			log.InfoContext(ctx, "songs not found")
			return nil, errs.NoSongs()
		}

		log.InfoContext(ctx, "get songs OK")

		return songs, nil
	}

	query := "SELECT id, song, group_name, release_date, text, link FROM songs WHERE true"
	argsCount := 1
	var args []any

	if fil.Song != "" {
		query += fmt.Sprintf(" AND LOWER(song)=LOWER($%d)", argsCount)
		args = append(args, fil.Song)
		argsCount++
	}

	if fil.Group != "" {
		query += fmt.Sprintf(" AND LOWER(group_name)=LOWER($%d)", argsCount)
		args = append(args, fil.Group)
		argsCount++
	}

	if fil.Date != nil {
		query += fmt.Sprintf(" AND release_date=$%d", argsCount)
		args = append(args, *fil.Date)
		argsCount++
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argsCount, argsCount+1)
	args = append(args, pag.Size, pag.Page)

	log.DebugContext(ctx, "get songs with filters", "args", args)

	songs, err := s.store.SongsByFilters(ctx, query, args)
	if err != nil {
		log.ErrorContext(ctx, "failed to get songs with filters", sl.Err(err))
		return nil, err
	}

	if len(songs) == 0 {
		log.InfoContext(ctx, "songs not found with filters", "fil", fil)
		return nil, errs.NoSongs()
	}

	log.InfoContext(ctx, "get songs with filters OK")

	return songs, nil
}

func (s *Service) GetSongText(ctx context.Context, pag types.Pagination, id int) ([]string, error) {
	ctx = logger.WithSongID(ctx, id)
	log := s.log.With(slog.String(fnName, getSongTextFn))

	log.DebugContext(ctx, "song text pagination", "pagination", pag)

	text, err := s.store.SongText(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, "failed to get song text", sl.Err(err))
		return nil, err
	}

	verses := strings.Split(text, "\n\n")

	log.DebugContext(ctx, "song verses", "verses", verses)

	if pag.Page <= 0 {
		return nil, errs.InvalidPage()
	}

	pageStart := (pag.Page - 1) * pag.Size
	pageEnd := pageStart + pag.Size

	if pageStart > len(verses)-1 {
		return nil, errs.EndOfText()
	}

	if pageEnd > len(verses) {
		pageEnd = len(verses)
	}

	if len(verses) == 0 {
		return nil, errs.NoText()
	}

	log.InfoContext(ctx, "get song text OK")

	return verses[pageStart:pageEnd], nil
}

func (s *Service) DeleteSong(ctx context.Context, id int) error {
	ctx = logger.WithSongID(ctx, id)
	log := s.log.With(slog.String(fnName, deleteSongFn))

	if err := s.store.DeleteSong(ctx, id); err != nil {
		log.ErrorContext(ctx, "failed to delete song", sl.Err(err))
		return err
	}

	log.InfoContext(ctx, "song delete OK")

	return nil
}

func (s *Service) UpdateSong(ctx context.Context, id int, req *types.UpdateSongRequest) error {
	ctx = logger.WithSongID(ctx, id)
	log := s.log.With(slog.String(fnName, updateSongFn))

	releaseDate, err := time.Parse(lib.Layout, req.ReleaseDate)
	if err != nil {
		log.ErrorContext(ctx, "failed to parse date", sl.Err(err))
		return errs.InvalidDate()
	}

	song := &types.Song{
		Song:        req.Song,
		Group:       req.Group,
		ReleaseDate: releaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}

	log.DebugContext(ctx, "to update", "req", req)

	if err := s.store.UpdateSong(ctx, id, song); err != nil {
		log.ErrorContext(ctx, "failed to update song", sl.Err(err))
	}

	log.InfoContext(ctx, "update song OK")

	return nil
}

func (s *Service) AddSong(ctx context.Context, req *types.SongRequest) error {
	log := s.log.With(slog.String(fnName, addSongFn))

	log.DebugContext(ctx, "song request", "req", req)

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	respch := make(chan Response)

	go func() {
		song, err := s.fetchSongDetails(ctx, req)
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
		log.WarnContext(ctx, "api call timeout")
		return errs.APICallTimeout()
	case resp := <-respch:
		if resp.err != nil {
			log.ErrorContext(ctx, "failed to fetch song details", sl.Err(resp.err))
			return resp.err
		}
		if err := s.store.AddSong(ctx, &resp.Song); err != nil {
			log.ErrorContext(ctx, "failed to add song", sl.Err(err))
			return err
		}
		return nil
	}
}

func (s *Service) fetchSongDetails(ctx context.Context, req *types.SongRequest) (*types.Song, error) {
	log := s.log.With(fnName, fetchSongDetailsFn)

	u, err := lib.ParseURL(s.url, req)
	if err != nil {
		log.ErrorContext(ctx, "failed to parse URL", sl.Err(err))
		return nil, err
	}

	log.InfoContext(ctx, "parsed URL", "url", u)

	client := &http.Client{}

	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(r)
	if err != nil {
		log.ErrorContext(ctx, "failed to make request to API", "url", u)
		return nil, err
	}
	defer resp.Body.Close()

	log.InfoContext(ctx, "response status code", "code", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.ErrorContext(ctx, "unexpected status code", "status code", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	details := new(types.Details)

	if err := json.NewDecoder(resp.Body).Decode(details); err != nil {
		log.ErrorContext(ctx, "failed to decode JSON response", sl.Err(err))
		return nil, err
	}

	releaseDate, err := time.Parse(lib.Layout, details.ReleaseDate)
	if err != nil {
		log.ErrorContext(ctx, "failed to parse date", sl.Err(err))
		return nil, err
	}

	song := &types.Song{
		Song:        req.Song,
		Group:       req.Group,
		ReleaseDate: releaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	log.InfoContext(ctx, "fetch song details OK", "details", details)

	return song, nil
}
