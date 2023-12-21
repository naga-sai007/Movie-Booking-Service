package security

import (
	"fmt"
	"strings"
	"theatre-service/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const (
	AdminUser    UserType = "admin"
	TheatreAdmin UserType = "theatre-admin"
	NormalUser   UserType = "user"
)

type UserType string

type UserClaims struct {
	UserId   int64
	UserName string
	UserType UserType
}

func Authenticate(c *fiber.Ctx) error {
	// Extract the raw token from the Authorization header
	rawToken := c.Get("Authorization")
	if rawToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Check for the custom prefix indicating an API key
	const apiKeyPrefix = "Apikey "
	if strings.HasPrefix(rawToken, apiKeyPrefix) {
		apiKey := strings.TrimPrefix(rawToken, apiKeyPrefix)

		// Validate the API key (replace this with your actual validation logic)
		if !isValidAPIKey(apiKey) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Set the user information in the context for subsequent handlers
		//user := getUserFromAPIKey(apiKey)
		//c.Locals("user", user)
	} else {
		// Assume it's a JWT with the "Bearer " prefix
		const jwtPrefix = "Bearer "
		tokenString := strings.TrimPrefix(rawToken, jwtPrefix)

		// Parse and validate the JWT
		claims := parseAndValidateJWT(tokenString)
		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Set the user information in the context for subsequent handlers
		user := getUserFromClaims(claims)
		c.Locals("user", user)
	}

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

func getUserFromClaims(claims jwt.MapClaims) UserClaims {
	// Extract user information from claims and create a User object
	user := UserClaims{
		UserId:   int64(claims["userId"].(float64)),
		UserName: claims["userName"].(string),
		UserType: UserType(claims["userType"].(string)),
	}
	return user
}

func Authorize(c *fiber.Ctx) error {
	// Extract the raw token from the Authorization header
	rawToken := c.Get("Authorization")

	// Check for the custom prefix indicating an API key
	const apiKeyPrefix = "Apikey "
	if strings.HasPrefix(rawToken, apiKeyPrefix) {
		// If it's an API key, skip authorization and proceed to the next handler
		return c.Next()
	}

	// If it's not an API key, assume it's a JWT
	user := c.Locals("user").(UserClaims)

	if user.UserType == AdminUser || user.UserType == TheatreAdmin {
		return c.Next()
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
}

func isValidAPIKey(apiKey string) bool {
	// Retrieve the securely stored hashed version of the API key
	storedAPIKey := utils.GetApiKeyFromConfig() // Fetch the hashed API key from your secure storage (e.g., database)

	// Use a secure hash comparison function (e.g., bcrypt.CompareHashAndPassword)
	err := bcrypt.CompareHashAndPassword([]byte(apiKey), []byte(storedAPIKey))
	return err == nil
}
