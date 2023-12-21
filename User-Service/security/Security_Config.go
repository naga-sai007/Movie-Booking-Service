package security

import (
	"fmt"
	"strings"
	"user-service/models"
	"user-service/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	// Extract JWT from headers
	tokenWithPrefix := c.Get("Authorization")
	if tokenWithPrefix == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Remove the "Bearer " prefix
	tokenString := strings.TrimPrefix(tokenWithPrefix, "Bearer ")

	// Parse and validate the token
	claims := parseAndValidateJWT(tokenString)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Set the user information in the context for subsequent handlers
	user := getUserFromClaims(claims)
	c.Locals("user", user)

	// Continue to the next handler
	return c.Next()
}

func parseAndValidateJWT(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.GetSecretKeyFromConfig()), nil
	})

	if err != nil || !token.Valid {
		fmt.Println("Error parsing or validating token:", err)
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Error extracting claims from token")
		return nil
	}

	//fmt.Println("Decoded Token Claims:", claims)
	return claims
}

func getUserFromClaims(claims jwt.MapClaims) models.User {
	// Extract user information from claims and create a User object
	user := models.User{
		UserID:   int64(claims["userId"].(float64)),
		UserName: claims["userName"].(string),
		UserType: models.UserType(claims["userType"].(string)),
	}
	return user
}

func AuthorizeAdmin(c *fiber.Ctx) error {
	// Retrieve the authenticated user from the context
	user := c.Locals("user").(models.User)

	// Check if the user has admin privileges
	if user.UserType != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	// Continue to the next handler
	return c.Next()
}

// both user and admin can access
func AuthorizeUser(c *fiber.Ctx) error {
	// Retrieve the authenticated user from the context
	user := c.Locals("user").(models.User)

	// Check if the user has the correct permissions to access the resource
	// Add your specific authorization logic based on user type or other criteria
	if user.UserType != models.UserType("theatre-admin") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	// Continue to the next handler
	return c.Next()
}
