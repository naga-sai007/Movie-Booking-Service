package handlers

import (
	"booking-service/models"
	"booking-service/repository"
	"booking-service/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateBooking(c *fiber.Ctx) error {

	var newBooking models.Booking
	if err := c.BodyParser(&newBooking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	var authToken string
	authTokenInterface := c.Locals("authToken")

	if authTokenInterface != nil {
		authToken, _ = authTokenInterface.(string)
	} else {
		// Handle the case where "authToken" is not found in Locals
		// This could be an error condition, depending on your requirements
		return fiber.NewError(fiber.StatusInternalServerError, "Auth token not found in Locals")
	}

	createdBooking, err := repository.CreateBooking(newBooking, authToken)
	if err != nil {
		log.Printf("Create booking error: %v", err)
		// Check if the error is a ValidationError
		if validationError, ok := err.(utils.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationError.Message})
		}

		// Handle other types of errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create booking"})
	}
	return c.Status(fiber.StatusCreated).JSON(createdBooking)

}

func GetBookingById(c *fiber.Ctx) error {
	var booking *models.Booking
	bookingId := c.Params("bookingId")
	booking, err := repository.GetBookingById(bookingId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Movie not found"})
	}
	return c.JSON(booking)
}

func BookingsByUser(c *fiber.Ctx) error {

	var bookings []models.Booking
	userId := c.Params("userId")
	converted_userID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	var authToken string
	authTokenInterface := c.Locals("authToken")

	if authTokenInterface != nil {
		authToken, _ = authTokenInterface.(string)
	} else {
		// Handle the case where "authToken" is not found in Locals
		// This could be an error condition, depending on your requirements
		return fiber.NewError(fiber.StatusInternalServerError, "Auth token not found in Locals")
	}

	bookings, err = repository.GetBookingsByUserId(converted_userID, authToken)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(bookings)
}

// func UpdateBooking(c *fiber.Ctx) error {
// 	return c.JSON(nil)
// }

// func DeleteBooking(c *fiber.Ctx) {}
