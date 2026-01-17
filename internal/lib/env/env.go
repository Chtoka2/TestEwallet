package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DotEnv struct{
	DbURL string
	Commission float64
	JWT_secret string
}

func Env_reader() DotEnv{
	if err := godotenv.Load(); err != nil{
		log.Fatal("File .env not exists")	
	}
	url_db := os.Getenv("dburl")
	jwt_secret:=os.Getenv("JWT_SECRET")
	return DotEnv{DbURL: url_db, JWT_secret: jwt_secret}
}