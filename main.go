package main

import (
	"flag"
	"log"
	"textscout/api"
	"textscout/common"
	"textscout/internal/populate"
)

func main() {
	// support two commands
	// 1. populate the database by parsing the json file
	// 2. start the REST server

	var commandFlag string
	var filePath string
	var searchBy string

	flag.StringVar(&commandFlag, "command", "", "which command to run. possible values are insertData and runServer")
	flag.StringVar(&filePath, "filePath", "", "path to the file to read from")
	flag.StringVar(&searchBy, "searchBy", "", "searchBy database or the inmemory inverted index. possible values are database and inmemIndex")
	flag.Parse()

	config := common.GetConfigOrDie()
	if commandFlag == "insertData" {
		if filePath == "" {
			log.Fatal("specify the filepath to read the data from.")
		}

		u := &populate.InsertData{
			FilePath: filePath,
			Config:   config,
		}
		u.InsertMovies()
	} else if commandFlag == "runServer" {
		if searchBy == "inmemIndex" && filePath == "" {
			log.Fatal("specify the filepath to read the data from.")
		}
		api.StartServer(config, searchBy, filePath)
	} else {
		log.Fatal("specify a valid command to run")
	}

}
