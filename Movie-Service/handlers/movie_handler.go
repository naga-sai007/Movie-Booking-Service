package handlers

import (
	"movie-service/models"
	"movie-service/repository"

	"github.com/gofiber/fiber/v2"
)

func GetMovies(c *fiber.Ctx) error {
	var movies []models.Movie = repository.GetAllMovies()
	return c.JSON(movies)
}

func GetMovieByID(c *fiber.Ctx) error {
	var movie *models.Movie
	movieId := c.Params("id")
	movie, err := repository.GetMovieByID(movieId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Movie not found"})
	}
	return c.JSON(movie)
}

func CreateMovie(c *fiber.Ctx) error {
	var newMovie models.Movie
	err := c.BodyParser(&newMovie)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	CreatedMovie, err := repository.CreateMovie(newMovie)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	return c.JSON(CreatedMovie)
}
