package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-service/models"
	"user-service/repository"
	"user-service/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User = repository.GetAllUsers()
	return c.JSON(users)

}

func GetUserByID(c *fiber.Ctx) error {
	var User *models.User
	UserId := c.Params("id")
	UserIdInt, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		return err
	}
	AuthUser := c.Locals("user").(models.User)
	if AuthUser.UserType == "user" && AuthUser.UserID != UserIdInt {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}
	User, err = repository.GetUserByID(UserId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(User)
}

func Signup(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return err
	}

	// Hash the user's password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.UserName = strings.ToUpper(user.UserName)
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	// Create the user in the database
	repository.CreateUser(*user)

	// Omit the password hash from the response
	user.Password = ""

	return c.JSON(user)
}

// Login handles user authentication and generates a JWT token
func Login(c *fiber.Ctx) error {

	var loginRequest models.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return err
	}

	user, err := repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}

	// Check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Password",
		})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.UserID
	claims["userName"] = user.UserName
	claims["userType"] = user.UserType
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (e.g., 24 hours)

	tokenString, err := token.SignedString([]byte(utils.GetSecretKeyFromConfig()))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": tokenString, "userId": user.UserID})
}

func RequestPasswordReset(c *fiber.Ctx) error {
	var PasswordChangeUser models.PasswordResetRequest

	// Parse JSON request body
	if err := c.BodyParser(&PasswordChangeUser); err != nil {
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse JSON"})
	}

	// Validate email and initiate the password reset process
	user, err := repository.FindUserByEmail(PasswordChangeUser.Email)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	token, err := utils.GenerateToken(15)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}
	expiry := time.Now().Add(time.Hour * 1) // Set expiration time (e.g., 1 hour)

	reset := models.PasswordReset{
		Email:  user.Email,
		UserID: user.UserID,
		Token:  token,
		Expiry: expiry,
	}

	err = repository.CreateResetItem(reset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create password reset record"})
	}

	// Send reset email with a link containing the token
	//sendResetEmail(user.Email, token)
	//this sendResetEmail method need to implement as of now directly giving the token to the user in json response body

	//fiber.Map{"message": "Password reset request initiated. Check your email for further instructions."} insted of token in Json
	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token, "email": user.Email})
}

// Function to handle the password reset process
func ResetPassword(c *fiber.Ctx) error {
	var ResetForm models.PasswordResetForm
	// Parse JSON request body
	if err := c.BodyParser(&ResetForm); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse JSON"})
	}

	// Validate the password reset token and process the reset
	PassReset, err := repository.ValidateResetToken(ResetForm)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Invalid or expired reset token"})
	}
	//return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Invalid or expired reset token"})
	//return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Reset token has expired"})
	// Update the user's password in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ResetForm.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	NewHashedPassword := string(hashedPassword)
	err = repository.UpdateUserPassword(PassReset.UserID, NewHashedPassword)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	// Delete the used password reset token from the database
	err = repository.DeleteResetItem(PassReset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete reset token"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Password reset successful. You can now log in with your new password."})
}
