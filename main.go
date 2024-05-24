package main

import (
	"circuithouse/api"
	"circuithouse/common"
	"circuithouse/internal/populate"
	"flag"
	"log"
)

func main() {
	// support two commands
	// 1. populate the database by parsing the json file
	// 2. start the REST server

	var commandFlag string
	var filePath string

	flag.StringVar(&commandFlag, "command", "", "which command to run. possible values are insertData and runServer")
	flag.StringVar(&filePath, "filePath", "", "path to the file to read from")
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
		api.StartServer(config)
	} else {
		log.Fatal("specify a valid command to run")
	}

}
