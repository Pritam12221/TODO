package main

import (
	db "TODO/database"
	s "TODO/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(".env file not found")
	}

	log.Println("host check", os.Getenv("DB_HOST"))
	log.Println("port check", os.Getenv("DB_PORT"))

	err = db.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		db.SSLMode(os.Getenv("DB_SSLMODE")),
	)
	if err != nil {
		log.Fatal(err)
	}

	srv := s.ServerRoutes()
	srv.Run()
}
