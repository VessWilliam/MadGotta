package main

import (
	"gotta/internal/routers"
	"log"
)

func main() {
	app, err := routers.NewApp()
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}

	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
