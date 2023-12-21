package handlers

import (
	"theatre-service/models"
	"theatre-service/repository"

	"github.com/gofiber/fiber/v2"
)

func GetShowByShowID(c *fiber.Ctx) error {
	var Show *models.Shows
	ShowId := c.Params("showid")
	Show, err := repository.GetShowByShowID(ShowId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Show not found"})
	}
	return c.JSON(Show)
}

func GetShowByMovieID(c *fiber.Ctx) error {
	var Shows []models.Shows
	ShowId := c.Params("movieid")
	Shows, err := repository.GetShowsByMovieID(ShowId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Show not found"})
	}
	return c.JSON(Shows)
}

func GetShowByTheatreID(c *fiber.Ctx) error {
	var Shows []models.Shows
	ShowId := c.Params("theatreid")
	Shows, err := repository.GetShowsByTheatreID(ShowId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Show not found"})
	}
	return c.JSON(Shows)
}

func CreateShow(c *fiber.Ctx) error {
	var newShow models.Shows
	err := c.BodyParser(&newShow)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	err = repository.CreateShow(newShow)
	if err != nil {
		// Handle the error and return it to the user
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create show", "details": err.Error()})
	}
	return c.JSON(newShow)
}

func UpdateShow(c *fiber.Ctx) error {
	ShowId := c.Params("showid")

	var updatedFields map[string]interface{}
	if err := c.BodyParser(&updatedFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	err := repository.UpdateTickets(ShowId, updatedFields)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	return c.SendStatus(fiber.StatusOK)
}
