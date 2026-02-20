package main

import (
	"asql/internal/server"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func run() error {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env not loaded")
	}
	server := server.NewServer()
	return server.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
