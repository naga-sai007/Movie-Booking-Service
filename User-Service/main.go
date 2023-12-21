package main

import (
	"fmt"
	databaseauth "user-service/database-auth"
	"user-service/handlers"
	"user-service/security"

	"github.com/gofiber/fiber/v2"
)

func main() {
	databaseauth.InitDB()
	db, err := databaseauth.DB.DB()

	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := fiber.New()

	app.Post("users/signup", handlers.Signup)

	app.Post("users/login", handlers.Login)

	// // Apply authentication middleware to all routes
	// Apply authorization middleware to /users route
	app.Get("/users", security.Authenticate, security.AuthorizeAdmin, handlers.GetUsers)

	// Apply authorization middleware to /users/:id route
	app.Get("/users/:id", security.Authenticate, handlers.GetUserByID)

	app.Post("/users/password-reset/request", handlers.RequestPasswordReset)
	app.Post("/users/password-reset/reset", handlers.ResetPassword)

	err = app.Listen(":3000")
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
