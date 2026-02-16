package main

import (
	"asql/internal/server"
	"log"
)

func run() error {
	server := server.NewServer()
	return server.Run("0.0.0.0:8080")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
