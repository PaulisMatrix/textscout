package common

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName string
	DBUser string
	DBPass string
}

func GetConfigOrDie() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("error in reading the .env file")
	}
	return &Config{
		DBName: os.Getenv("POSTGRES_DB"),
		DBUser: os.Getenv("POSTGRES_USER"),
		DBPass: os.Getenv("POSTGRES_PASS"),
	}
}
