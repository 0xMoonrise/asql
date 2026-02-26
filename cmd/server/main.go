package main

import (
	"asql/internal/server"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var HOST string = "127.0.0.1"
var PORT string = "8080"

func run() error {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env not loaded")
	}

	server := server.NewServer()

	h := os.Getenv("HOST")
	if h != "" {
		HOST = h
	}

	p := os.Getenv("PORT")
	if p != "" {
		PORT = p
	}

	return server.Run(HOST + ":" + PORT)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
