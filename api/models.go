package api

type Results struct {
	Results []MovieData `json:"results"`
}

// search by title and overview(desc)
type MovieData struct {
	Adult         bool    `json:"adult,omitempty"`
	BackdropPath  string  `json:"backdrop_path,omitempty"`
	GenreIDs      []int32 `json:"genre_ids,omitempty"`
	ID            int32   `json:"id"`
	Language      string  `json:"original_language,omitempty"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview,omitempty"`
	Popularity    float64 `json:"popularity,omitempty"`
	PosterPath    string  `json:"poster_path,omitempty"`
	ReleaseDate   string  `json:"release_date,omitempty"`
	MovieTitle    string  `json:"title"`
	Video         bool    `json:"video,omitempty"`
	VoteAverage   float64 `json:"vote_average,omitempty"`
	VoteCount     int64   `json:"vote_count,omitempty"`
}
