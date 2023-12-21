package repository

import (
	"errors"
	"strconv"
	"time"
	databaseauth "user-service/database-auth"
	"user-service/models"
)

func CreateResetItem(Reset models.PasswordReset) error {

	Reset.CreatedAt = time.Now()
	if err := databaseauth.DB.Create(&Reset).Error; err != nil {
		return err
	}
	return nil
}

func ValidateResetToken(Reset models.PasswordResetForm) (models.PasswordReset, error) {
	var PassReset models.PasswordReset
	result := databaseauth.DB.Where("user_id = ?", Reset.UserId).Order("created_at desc").First(&PassReset)
	if result.Error != nil {
		return PassReset, errors.New("reset ")
	}
	if PassReset.Token != Reset.Token {
		return PassReset, errors.New("token err")
	}
	if time.Now().After(PassReset.Expiry) {
		return PassReset, errors.New("token Expired")
	}

	return PassReset, nil
}

func UpdateUserPassword(ResetUserId int64, NewHashedPassword string) error {

	userId := strconv.FormatInt(ResetUserId, 10)
	var existingUser *models.User
	existingUser, err := GetUserByID(userId)
	if err != nil {
		return errors.New("failed to retrieve User")
	}
	existingUser.Password = NewHashedPassword
	result := databaseauth.DB.Save(existingUser)
	if result.Error != nil {
		return errors.New("failed to save User")
	}
	return nil
}

func DeleteResetItem(reset models.PasswordReset) error {
	result := databaseauth.DB.Delete(reset)
	if result.Error != nil {
		return errors.New("unable to delete Reset Item")
	}
	return nil
}
