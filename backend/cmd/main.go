package main

import (
	"chatbot-server/internal/app"
	"log"
)

func main() {
	application, err := app.InitializeApp()

	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
