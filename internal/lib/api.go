package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/erknas/song-library/internal/types"
)

const Layout = "02.01.2006"

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

func MakeHTTPFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
		defer cancel()

		if err := fn(ctx, w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ParseID(r *http.Request) (int, error) {
	id := r.FormValue("id")
	return strconv.Atoi(id)
}

func PaginationValues(r *http.Request) (types.Pagination, error) {
	var (
		strPage = r.FormValue("page")
		strSize = r.FormValue("size")
	)

	if len(strPage) == 0 {
		strPage = "1"
	}

	if strSize != "10" && strSize != "25" && strSize != "50" {
		strSize = "10"
	}

	page, err := strconv.Atoi(strPage)
	if err != nil {
		return types.Pagination{}, err
	}

	size, err := strconv.Atoi(strSize)
	if err != nil {
		return types.Pagination{}, err
	}

	pagination := types.Pagination{
		Page: page,
		Size: size,
	}

	return pagination, nil
}

func FilterValues(r *http.Request) (types.Filter, error) {
	var (
		song    = r.FormValue("song")
		group   = r.FormValue("group")
		strDate = r.FormValue("date")
	)

	if len(strDate) == 0 {
		filters := types.Filter{
			Song:  song,
			Group: group,
			Date:  nil,
		}
		return filters, nil
	}

	date, err := time.Parse(Layout, strDate)
	if err != nil {
		return types.Filter{}, err
	}

	filters := types.Filter{
		Song:  song,
		Group: group,
		Date:  &date,
	}

	return filters, nil
}
