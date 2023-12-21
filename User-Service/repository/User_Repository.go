package repository

import (
	"errors"
	databaseauth "user-service/database-auth"
	"user-service/models"
)

func CreateUser(newUser models.User) {

	databaseauth.DB.Create(&newUser)

}

func FindUserByEmail(email string) (models.User, error) {
	var user models.User
	result := databaseauth.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	if !user.IsActive {
		return models.User{}, errors.New("user not active")
	}
	return user, nil
}

func GetAllUsers() []models.User {
	var Users []models.User
	databaseauth.DB.Find(&Users)
	return Users
}

func GetUserByID(UserID string) (*models.User, error) {
	var User models.User
	if err := databaseauth.DB.First(&User, UserID).Error; err != nil {
		return nil, err
	}
	return &User, nil
}

func DeleteUser(UserID string) {

	databaseauth.DB.Delete(&models.User{}, UserID)
}
