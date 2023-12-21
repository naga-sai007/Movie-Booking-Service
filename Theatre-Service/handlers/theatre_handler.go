package handlers

import (
	"theatre-service/models"
	"theatre-service/repository"

	"github.com/gofiber/fiber/v2"
)

func GetAllTheatres(c *fiber.Ctx) error {
	var Theatres []models.Theatre = repository.GetAllTheatres()
	return c.JSON(Theatres)
}

func GetTheatreByID(c *fiber.Ctx) error {
	var Theatre *models.Theatre
	TheatreId := c.Params("theatreid")
	Theatre, err := repository.GetTheatreByID(TheatreId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Theatre not found"})
	}
	return c.JSON(Theatre)
}

func CreateTheatre(c *fiber.Ctx) error {
	var newTheatre models.Theatre
	err := c.BodyParser(&newTheatre)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	repository.CreateTheatre(newTheatre)
	return c.JSON(newTheatre)
}
