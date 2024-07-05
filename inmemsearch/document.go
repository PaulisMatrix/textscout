package inmemsearch

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"textscout/common"
)

type Document struct {
	ID            int
	Adult         bool
	BackdropPath  string
	GenreIDs      []int32
	MovieID       int32
	Language      string
	OriginalTitle string
	Overview      string
	Popularity    float64
	PosterPath    string
	ReleaseDate   string
	MovieTitle    string
	Video         bool
	VoteAverage   float64
	VoteCount     int64
}

func loadMovies(filePath string) ([]Document, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		log.Fatal("failed to open the file", err)
		return []Document{}, err
	}
	defer fd.Close()

	// start reading the json file in-memory
	jsonBytes, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal("failed to read the file", err)
		return []Document{}, err
	}

	var results common.Results
	err = json.Unmarshal(jsonBytes, &results)
	if err != nil {
		log.Fatal("failed to unmarshal the json", err)
	}
	return marshalToDocs(results.Results), nil
}

func marshalToDocs(movieData []common.MovieData) []Document {
	docs := make([]Document, 0)

	for idx, md := range movieData {
		d := Document{
			ID:            idx,
			Adult:         md.Adult,
			BackdropPath:  md.BackdropPath,
			GenreIDs:      md.GenreIDs,
			MovieID:       md.ID,
			Language:      md.Language,
			OriginalTitle: md.OriginalTitle,
			Overview:      md.Overview,
			Popularity:    md.Popularity,
			PosterPath:    md.PosterPath,
			ReleaseDate:   md.ReleaseDate,
			MovieTitle:    md.MovieTitle,
			Video:         md.Video,
			VoteAverage:   md.VoteAverage,
			VoteCount:     md.VoteCount,
		}
		docs = append(docs, d)
	}
	return docs
}
