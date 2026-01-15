package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DotEnv struct{
	DbURL string
}

func Env_reader() DotEnv{
	if err := godotenv.Load(); err != nil{
		log.Fatal("File .env not exists")	
	}
	url_db := os.Getenv("dburl")
	return DotEnv{DbURL: url_db}
}