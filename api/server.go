package api

import (
	"circuithouse/common"
	"circuithouse/internal/database"
	"fmt"
	"net/http"
)

type SearchAPI struct {
	querier database.Querier
}

func (s *SearchAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
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

	http.Handle("/api/v1/search", s)

	fmt.Println("starting the server at port 8080")
	err = http.ListenAndServe(":8080", s)
	if err != nil {
		panic(err.Error())
	}

}
