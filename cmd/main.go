package main

import (
	"log"
	"os"

	"identify/internal/app"
	"identify/internal/database"
	"identify/internal/di"
)

func main() {
	// Build dependency container
	container, err := di.BuildContainer()
	if err != nil {
		log.Fatal(err)
	}

	// Ensure DB is closed when application shuts down
	defer database.DisconnectPostgresDb(container.DB)

	// Create cancellable context listening for SIGINT/SIGTERM
	ctx, cancel := app.WithContext()
	defer cancel()

	// Instantiate and run the App
	a := app.NewApp(container.Registry, container.Logger)
	if err := a.Run(ctx); err != nil {
		container.Logger.Error("Application exited with error", "error", err)
		os.Exit(1)
	}
}
