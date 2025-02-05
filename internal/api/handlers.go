package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/erknas/song-library/internal/errs"
	"github.com/erknas/song-library/internal/lib"
	"github.com/erknas/song-library/internal/types"
)

func (s *Server) handleSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetSongText(ctx, w, r)
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

//	@Summary		Get a list of songs
//	@Description	Get a paginated list of songs with optional filtering
//	@Tags			songs
//	@Produce		json
//	@Param			page	query		int		false	"Page number"				default(1)	example(1)
//	@Param			size	query		int		false	"Number of songs per page"	default(10)	example(10)	Enums(10,25,50)
//	@Param			song	query		string	false	"Filter by song"			example(Supermassive Black Hole)
//	@Param			group	query		string	false	"Filter by group"			example(Muse)
//	@Param			date	query		string	false	"Filter by release_date"	example(16.07.2006)
//	@Success		200		{object}	[]types.Song
//	@Failure		400		{object}	errs.APIError
//	@Failure		500		{string}	internal	server	error
//	@Router			/songs [get]
func (s *Server) handleGetSongs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	pag, err := lib.SongsPaginationValues(r)
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

//	@Summary		Get a list of verses by song
//	@Description	Get a paginated list of verses by song
//	@Tags			song
//	@Produce		json
//	@Param			id		query		int	true	"Song ID"
//	@Param			page	query		int	false	"Page number"				default(1)	example(1)
//	@Param			size	query		int	false	"Number of verses per page"	default(1)	example(1)	Enums(1,5,10)
//	@Success		200		{object}	[]types.Text
//	@Failure		400		{object}	errs.APIError
//	@Failure		500		{string}	internal	server	error
//	@Router			/song [get]
func (s *Server) handleGetSongText(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := lib.ParseID(r)
	if err != nil {
		return errs.InvalidID()
	}

	pag, err := lib.TextPaginationValues(r)
	if err != nil {
		return err
	}

	text, err := s.srv.GetSongText(ctx, pag, id)
	if err != nil {
		return err
	}

	resp := types.Text{Text: text}

	return lib.WriteJSON(w, http.StatusOK, resp)
}

//	@Summary		Delete song
//	@Description	Delete song by ID
//	@Tags			song
//	@Produce		json
//	@Param			id	query		int	true	"Song ID"
//	@Success		200	{object}	types.SongResponse
//	@Failure		400	{object}	errs.APIError
//	@Failure		500	{string}	internal	server	error
//	@Router			/song [delete]
func (s *Server) handleDeleteSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := lib.ParseID(r)
	if err != nil {
		return errs.InvalidID()
	}

	if err := s.srv.DeleteSong(ctx, id); err != nil {
		return err
	}

	resp := types.NewSongResponse(http.StatusOK, "song successfully deleted")

	return lib.WriteJSON(w, http.StatusOK, resp)
}

//	@Summary		Update song
//	@Description	Update song by ID
//	@Tags			song
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int						true	"Song ID"
//	@Param			song	body		types.UpdateSongRequest	true	"Update song data"
//	@Success		200		{object}	[]types.SongResponse
//	@Failure		400		{object}	errs.APIError
//	@Failure		500		{string}	internal	server	error
//	@Router			/song [put]
func (s *Server) handleUpdateSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := lib.ParseID(r)
	if err != nil {
		return errs.InvalidID()
	}

	req := new(types.UpdateSongRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return errs.InvalidJSON()
	}
	defer r.Body.Close()

	if err := s.srv.UpdateSong(ctx, id, req); err != nil {
		return err
	}

	resp := types.NewSongResponse(http.StatusOK, "song successfully updated")

	return lib.WriteJSON(w, http.StatusOK, resp)
}

//	@Summary		Add song
//	@Description	Add song
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			song	body		types.SongRequest	true	"Song data"
//	@Failure		400		{object}	errs.APIError
//	@Failure		500		{string}	internal	server	error
//	@Router			/songs [post]
func (s *Server) handleAddSong(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := new(types.SongRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return errs.InvalidJSON()
	}
	defer r.Body.Close()

	if err := s.srv.AddSong(ctx, req); err != nil {
		return err
	}

	resp := types.NewSongResponse(http.StatusOK, "song successfully added")

	return lib.WriteJSON(w, http.StatusOK, resp)
}
