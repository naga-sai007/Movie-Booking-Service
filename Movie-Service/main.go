package main

import (
	"fmt"
	databaseauth "movie-service/database-auth"
	"movie-service/handlers"
	"movie-service/security"

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
	app.Get("/movies", handlers.GetMovies)
	app.Get("/movies/:id", handlers.GetMovieByID)
	app.Post("/movies", security.Authenticate, security.AuthorizeAdmin, handlers.CreateMovie) //authorization and authentication needed (admin or theatre admin only can add)

	err = app.Listen(":3001")
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
