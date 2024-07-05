package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"textscout/common"
	textsearch "textscout/inmemsearch"
	"textscout/internal/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SearchAPI struct {
	querier       database.Querier
	inMemoryIndex *textsearch.InMemSearch
	searchBy      string
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

		log.Printf("Method: %s and URI: %s\n", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)

	})
	return f
}

func (s *SearchAPI) readyResponseDB(resp []database.Movie) common.Response {
	movies := []common.MovieData{}

	for _, data := range resp {
		d := common.MovieData{
			Adult:         data.Adult.Bool,
			BackdropPath:  data.BackdropPath.String,
			GenreIDs:      data.GenreIds,
			ID:            data.MovieID,
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

	return common.Response{
		Movies: movies,
	}

}

func (s *SearchAPI) readyResponseInMemIndex(resp []textsearch.Document) common.Response {
	movies := []common.MovieData{}

	for _, data := range resp {
		d := common.MovieData{
			Adult:         data.Adult,
			BackdropPath:  data.BackdropPath,
			GenreIDs:      data.GenreIDs,
			ID:            data.MovieID,
			Language:      data.Language,
			OriginalTitle: data.OriginalTitle,
			Overview:      data.Overview,
			Popularity:    data.Popularity,
			PosterPath:    data.PosterPath,
			ReleaseDate:   data.ReleaseDate,
			MovieTitle:    data.MovieTitle,
			Video:         data.Video,
			VoteAverage:   data.VoteAverage,
			VoteCount:     data.VoteCount,
		}
		movies = append(movies, d)
	}

	return common.Response{
		Movies: movies,
	}

}

func (s *SearchAPI) validateAndWriteAPIResponseDatabase(w http.ResponseWriter, dbResp []database.Movie) {
	if len(dbResp) == 0 {
		http.Error(w, "no records found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(s.readyResponseDB(dbResp))
	if err != nil {
		log.Fatalf("failed to marshal the resp: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (s *SearchAPI) validateAndWriteAPIResponseInMemIndex(w http.ResponseWriter, inMemResp []textsearch.Document) {
	if len(inMemResp) == 0 {
		http.Error(w, "no records found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(s.readyResponseInMemIndex(inMemResp))
	if err != nil {
		log.Fatalf("failed to marshal the resp: %+v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (s *SearchAPI) useDatabase(w http.ResponseWriter, title string, desc string) {

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
		s.validateAndWriteAPIResponseDatabase(w, resp)
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
		s.validateAndWriteAPIResponseDatabase(w, resp)
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
		s.validateAndWriteAPIResponseDatabase(w, resp)
	}
}

func (s *SearchAPI) useInMemoryIndex(w http.ResponseWriter, title string, desc string) {
	matchedDocIds := s.inMemoryIndex.Intersection(common.ConcatStrings(title, desc))

	s.validateAndWriteAPIResponseInMemIndex(w, matchedDocIds)

}

func (s *SearchAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// search the db given the search query/queries
	values := r.URL.Query()
	title := values.Get("title")
	desc := values.Get("desc")

	if s.searchBy == "inmemIndex" {
		s.useInMemoryIndex(w, title, desc)
	} else {
		s.useDatabase(w, title, desc)
	}

}

func initDB(dbName, dbUser, dbPass string) *database.Queries {
	pgDB, err := database.NewPostgres(dbName, dbUser, dbPass)
	if err != nil {
		panic(err.Error())
	}
	queries := database.New(pgDB.DB)
	return queries
}

func getHandler(config *common.Config, searchBy string, filePath string) *SearchAPI {
	if searchBy == "inmemIndex" {
		log.Println("using the in-memory index for searching")
		return &SearchAPI{
			inMemoryIndex: textsearch.GetInMemSearch(filePath),
			searchBy:      searchBy,
		}
	} else {
		log.Println("using the database for searching")
		return &SearchAPI{
			querier:  initDB(config.DBName, config.DBUser, config.DBPass),
			searchBy: searchBy,
		}
	}
}

func StartServer(config *common.Config, searchBy string, filePath string) {
	// REST server
	// One Endpoint: localhost:8080/api/v1/search?title=""&desc=""

	s := getHandler(config, searchBy, filePath)

	handler := Validator(Logger(s))
	http.Handle("/api/v1/search", handler)

	log.Println("starting the server at port 8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err.Error())
	}

}
