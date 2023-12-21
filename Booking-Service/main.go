package main

import (
	databaseauth "booking-service/database-auth"
	"booking-service/handlers"
	"booking-service/security"
	"fmt"

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

	app.Get("/bookings/:bookingId", security.Authenticate, handlers.GetBookingById)   //get booking by bookingId
	app.Get("/user/bookings/:userId", security.Authenticate, handlers.BookingsByUser) //get all bookings of User by UserID

	app.Post("/bookings", security.Authenticate, handlers.CreateBooking) //createbooking
	//app.Put("/bookings/:bookingId",handlers)    //update booking
	//app.Delete("/bookings/:bookingId",handlers) //deletebooking

	err = app.Listen(":3003")
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}

}
