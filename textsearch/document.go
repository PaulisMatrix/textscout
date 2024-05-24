package textsearch

import (
	"circuithouse/api"
	"encoding/json"
	"io"
	"log"
	"os"
)

func loadMovies(filePath string) ([]api.MovieData, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		log.Fatal("failed to open the file", err)
		return []api.MovieData{}, err
	}
	defer fd.Close()

	// start reading the json file in-memory
	jsonBytes, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal("failed to read the file", err)
		return []api.MovieData{}, err
	}

	var results api.Results
	err = json.Unmarshal(jsonBytes, &results)
	if err != nil {
		log.Fatal("failed to unmarshal the json", err)
	}
	return results.Results, nil
}
