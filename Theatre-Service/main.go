package main

import (
	"log"
	databaseauth "theatre-service/database-auth"
	"theatre-service/handlers"
	"theatre-service/security"

	"github.com/gofiber/fiber/v2"
)

func main() {

	databaseauth.InitDb()
	db, err := databaseauth.DB.DB()

	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := fiber.New()

	app.Get("/theatres", handlers.GetAllTheatres)
	app.Get("/theatres/:theatreid", handlers.GetTheatreByID)
	app.Post("/theatres", security.Authenticate, security.Authorize, handlers.CreateTheatre) // Create a theatre (authorize the user as theatre admin)

	app.Post("/theatres/:theatreid/shows", security.Authenticate, security.Authorize, handlers.CreateShow) // Create a show for a specific theatre (authenticate theatre admin)

	app.Get("/movies/:movieid/shows", handlers.GetShowByMovieID)       // List of shows for a specific movie
	app.Get("/theatres/:theatreid/shows", handlers.GetShowByTheatreID) // List of shows for a specific theatre
	app.Get("/shows/:showid", handlers.GetShowByShowID)                // Get show information for booking purposes

	app.Put("/shows/tickets/:showid", security.Authenticate, security.Authorize, handlers.UpdateShow)

	err = app.Listen(":3002")
	if err != nil {
		log.Printf("Error starting the server: %f", err)
	}
}
