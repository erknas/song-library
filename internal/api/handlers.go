package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/erknas/song-library/internal/lib"
	"github.com/erknas/song-library/internal/types"
)

func (s *Server) handleSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetSongLyrics(ctx, w, r)
	case http.MethodDelete:
		return s.handleDeleteSong(ctx, w, r)
	case http.MethodPut:
		return s.handleUpdateSong(ctx, w, r)
	default:
		return lib.WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}
}

func (s *Server) handleSongs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetSongs(ctx, w, r)
	case http.MethodPost:
		return s.handleAddSong(ctx, w, r)
	default:
		return lib.WriteJSON(w, http.StatusMethodNotAllowed, nil)
	}
}

func (s *Server) handleGetSongs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	pag, err := lib.PaginationValues(r)
	if err != nil {
		return err
	}

	fil, err := lib.FilterValues(r)
	if err != nil {
		return err
	}

	songs, err := s.srv.GetSongs(ctx, pag, fil)
	if err != nil {
		return err
	}

	return lib.WriteJSON(w, http.StatusOK, types.Songs{Songs: songs})
}

func (s *Server) handleGetSongLyrics(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := lib.ParseID(r)
	if err != nil {
		return err
	}

	pag, err := lib.PaginationValues(r)
	if err != nil {
		return err
	}

	text, err := s.srv.GetSongText(ctx, pag, id)
	if err != nil {
		return err
	}

	return lib.WriteJSON(w, http.StatusOK, text)
}

func (s *Server) handleDeleteSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := lib.ParseID(r)
	if err != nil {
		return err
	}

	if err := s.srv.DeleteSong(ctx, id); err != nil {
		return err
	}

	resp := types.NewSongResponse(http.StatusOK, "song successfully deleted")

	return lib.WriteJSON(w, http.StatusOK, resp)
}

func (s *Server) handleUpdateSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := lib.ParseID(r)
	if err != nil {
		return err
	}

	req := new(types.Song)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	defer r.Body.Close()

	if err := s.srv.UpdateSong(ctx, id, req); err != nil {
		return err
	}

	resp := types.NewSongResponse(http.StatusOK, "song successfully updated")

	return lib.WriteJSON(w, http.StatusOK, resp)
}

func (s *Server) handleAddSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := new(types.SongRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	defer r.Body.Close()

	if err := s.srv.AddSong(ctx, req); err != nil {
		s.log.Error("failed to add song", "error", err)
		return err
	}

	resp := types.NewSongResponse(http.StatusOK, "song successfully added")

	return lib.WriteJSON(w, http.StatusOK, resp)
}
