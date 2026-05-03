package main

import (
	"log"
	"os"

	"cultivation-client/internal/app"
)

func main() {
	application := app.New()

	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
		os.Exit(1)
	}
}
