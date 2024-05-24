package api

import (
	"circuithouse/common"
	"circuithouse/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SearchAPI struct {
	querier database.Querier
}

// Middlewares:
// 1. Validator: validate its a GET request, check for at least one query params present to search.
// 2. Logger: log every incoming request

func Validator(next http.Handler) http.Handler {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		queryValues := r.URL.Query()
		title := queryValues.Get("title")
		desc := queryValues.Get("desc")
		if title == "" && desc == "" {
			http.Error(w, "At least one query parameter (title or desc) is required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
	return f
}

func Logger(next http.Handler) http.Handler {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Method: %s and URI: %s\n", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)

	})
	return f
}

func (s *SearchAPI) readyResponse(dbResp []database.Movie) Response {
	movies := []MovieData{}

	for _, data := range dbResp {
		d := MovieData{
			Adult:         data.Adult.Bool,
			BackdropPath:  data.BackdropPath.String,
			GenreIDs:      data.GenreIds,
			ID:            data.ID,
			Language:      data.MovieLanguage.String,
			OriginalTitle: data.MovieOriginalTitle.String,
			Overview:      data.MovieOverview.String,
			Popularity:    data.Popularity.Float64,
			PosterPath:    data.PosterPath.String,
			ReleaseDate:   data.ReleaseDate.String,
			MovieTitle:    data.MovieTitle,
			Video:         data.Video.Bool,
			VoteAverage:   data.VoteAverage.Float64,
			VoteCount:     data.VoteCount.Int64,
		}
		movies = append(movies, d)
	}

	return Response{
		Movies: movies,
	}

}

func (s *SearchAPI) validateAndWriteAPIResponse(w http.ResponseWriter, dbResp []database.Movie) {
	if len(dbResp) == 0 {
		http.Error(w, "no records found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(s.readyResponse(dbResp))
	if err != nil {
		log.Fatalf("failed to marshal the resp: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (s *SearchAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// search the db given the search query/queries
	values := r.URL.Query()
	title := values.Get("title")
	desc := values.Get("desc")

	if title != "" && desc != "" {
		// search by both and return the results
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		args := database.GetMovieByTitleAndDescParams{
			Column1: pgtype.Text{
				String: title,
				Valid:  true,
			},
			Column2: pgtype.Text{
				String: desc,
				Valid:  true,
			},
		}

		resp, err := s.querier.GetMovieByTitleAndDesc(context, args)
		if err != nil {
			log.Fatalf("error while searching by both title and desc: %+v", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		s.validateAndWriteAPIResponse(w, resp)
	} else if title != "" {
		// search by the title
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		args := pgtype.Text{
			String: title,
			Valid:  true,
		}

		resp, err := s.querier.GetMovieByTitle(context, args)
		if err != nil {
			log.Fatalf("error while searching by title: %+v", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		s.validateAndWriteAPIResponse(w, resp)
	} else {
		// search by the desc
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		args := pgtype.Text{
			String: desc,
			Valid:  true,
		}

		resp, err := s.querier.GetMovieByDesc(context, args)
		if err != nil {
			log.Fatalf("error while searching by title: %+v", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		s.validateAndWriteAPIResponse(w, resp)
	}

}

func StartServer(config *common.Config) {
	// REST server
	// One Endpoint: localhost:8080/api/v1/search?title=""&desc=""
	pgDB, err := database.NewPostgres(config.DBName, config.DBUser, config.DBPass)
	if err != nil {
		panic(err.Error())
	}

	queries := database.New(pgDB.DB)
	s := &SearchAPI{
		querier: queries,
	}

	handler := Validator(Logger(s))
	http.Handle("/api/v1/search", handler)

	fmt.Println("starting the server at port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err.Error())
	}

}
