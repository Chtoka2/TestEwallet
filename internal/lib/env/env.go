package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Env_reader() string{
	if err := godotenv.Load(); err != nil{
		log.Fatal("File .env not exists")	
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == ""{
		log.Fatal("dbUrl is empty")
	}
	return dbUrl
}