package security

import (
	"booking-service/utils"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
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
	c.Locals("authToken", tokenWithPrefix)
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

func AuthorizeAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(UserClaims)

	if user.UserType == AdminUser || user.UserType == TheatreAdmin {
		return c.Next()
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
}
