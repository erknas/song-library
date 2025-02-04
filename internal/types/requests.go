package types

type SongRequest struct {
	Song  string `json:"song"`
	Group string `json:"group"`
}

type UpdateSongRequest struct {
	Song        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
