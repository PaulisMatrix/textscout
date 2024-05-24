package populate

import (
	"circuithouse/api"
	"circuithouse/common"
	"circuithouse/internal/database"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
)

// reads the json from the filepath and inserts into postgres
type InsertData struct {
	FilePath string
	Config   *common.Config
}

func (d *InsertData) InsertMovies() {
	// initialise the db connection
	postgres, err := database.NewPostgres(d.Config.DBName, d.Config.DBUser, d.Config.DBPass)
	if err != nil {
		panic(err.Error())
	}

	moviesDB := database.New(postgres.DB)

	fd, err := os.Open(d.FilePath)
	if err != nil {
		log.Fatal("failed to open the file", err)
	}
	defer fd.Close()

	// start reading the json file in-memory
	jsonBytes, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal("failed to read the file", err)
	}

	var results api.Results
	err = json.Unmarshal(jsonBytes, &results)
	if err != nil {
		log.Fatal("failed to unmarshal the json", err)
	}

	for _, movieData := range results.Results {
		arg := database.AddMovieParams{
			Adult: pgtype.Bool{
				Bool:  movieData.Adult,
				Valid: true,
			},
			BackdropPath: pgtype.Text{
				String: movieData.BackdropPath,
				Valid:  true,
			},
			GenreIds: movieData.GenreIDs,
			MovieID:  movieData.ID,
			MovieLanguage: pgtype.Text{
				String: movieData.Language,
				Valid:  true,
			},
			MovieOriginalTitle: pgtype.Text{
				String: movieData.OriginalTitle,
				Valid:  true,
			},
			MovieOverview: pgtype.Text{
				String: movieData.Overview,
				Valid:  true,
			},
			Popularity: pgtype.Float8{
				Float64: movieData.Popularity,
				Valid:   true,
			},
			PosterPath: pgtype.Text{
				String: movieData.PosterPath,
				Valid:  true,
			},
			ReleaseDate: pgtype.Text{
				String: movieData.ReleaseDate,
				Valid:  true,
			},
			MovieTitle: movieData.MovieTitle,
			Video: pgtype.Bool{
				Bool:  movieData.Video,
				Valid: true,
			},
			VoteAverage: pgtype.Float8{
				Float64: movieData.VoteAverage,
				Valid:   true,
			},
			VoteCount: pgtype.Int8{
				Int64: movieData.VoteCount,
				Valid: true,
			},
		}
		err := moviesDB.AddMovie(context.Background(), arg)
		if err != nil {
			log.Fatalf("error inserting movie data: %+v", err)
		}

	}

}
