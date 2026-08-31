package main

import (
	"log"

	"github.com/joho/godotenv"

	"jakes-resume-compiler/server"
	"jakes-resume-compiler/server/shared/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env vars")
	}

	app := config.NewApp()

	server.Start(app)
}
