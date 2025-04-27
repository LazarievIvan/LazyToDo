package main

import (
	"LazyToDo/internal/server"
	"log"
	"os"
)

// Entry point: starts the server on port defined in environmental variables.
func main() {
	port := os.Getenv("PORT")
	if err := server.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
