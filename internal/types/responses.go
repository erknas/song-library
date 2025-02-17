package types

type Text struct {
	Text []string `json:"text"`
}

type SongResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func NewSongResponse(status int, msg string) SongResponse {
	return SongResponse{
		StatusCode: status,
		Msg:        msg,
	}
}
