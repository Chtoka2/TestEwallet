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
	host := os.Getenv("host")
	if host == ""{
		log.Fatal("Host is empty")
	}
	user := os.Getenv("user")
	if user == ""{
		log.Fatal("User is empty")
	}
	dbname := os.Getenv("dbname")
	if dbname == ""{
		log.Fatal("Dbname is empty")
	}
	password := os.Getenv("password")
	if password == ""{
		log.Fatal("Password is empty")
	}
	sslmode := os.Getenv("sslmode")
	if sslmode == ""{
		sslmode="disable"
	}
	return "host="+host+" user="+user+" dbname="+dbname+" password="+password+" sslmode="+sslmode
}