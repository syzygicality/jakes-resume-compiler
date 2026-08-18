package main

import (
	"log"

	"github.com/joho/godotenv"

	"jakes-resume-compiler/server/config"
	"jakes-resume-compiler/server/grpc"
	"jakes-resume-compiler/server/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env vars")
	}

	app := config.NewApp()

	if app.Settings.Mode == "HTTP" {
		http.Start(app)
	} else {
		grpc.Start(app)
	}
}
