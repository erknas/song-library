package types

import "time"

type Song struct {
	ID          int       `json:"id"`
	Song        string    `json:"song"`
	Group       string    `json:"group"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type Details struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type Songs struct {
	Songs []*Song `json:"songs"`
}

type Pagination struct {
	Page int
	Size int
}

type Filter struct {
	Song  string
	Group string
	Date  *time.Time
}
