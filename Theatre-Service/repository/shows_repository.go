package repository

import (
	"errors"
	databaseauth "theatre-service/database-auth"
	"theatre-service/models"
)

func CreateShow(newShow models.Shows) error {
	result := databaseauth.DB.Create(&newShow)
	if result.Error != nil {
		// Handle the error and return it
		return result.Error
	}
	return nil
}

func GetShowByShowID(showID string) (*models.Shows, error) {
	var Show models.Shows
	if err := databaseauth.DB.First(&Show, showID).Error; err != nil {
		return nil, err
	}
	return &Show, nil
}

func GetShowsByMovieID(movieID string) ([]models.Shows, error) {
	var Shows []models.Shows
	if err := databaseauth.DB.Where("movie_id = ?", movieID).Find(&Shows).Error; err != nil {
		return nil, err
	}
	return Shows, nil
}

func GetShowsByTheatreID(theatreID string) ([]models.Shows, error) {
	var Shows []models.Shows
	if err := databaseauth.DB.Where("theatre_id = ?", theatreID).Find(&Shows).Error; err != nil {
		return nil, err
	}
	return Shows, nil
}

func UpdateTickets(showId string, updatedFields map[string]interface{}) error {

	existingShow, err := GetShowByShowID(showId)
	if err != nil {
		return err
	}
	if !existingShow.ShowStatus {
		return errors.New("show status is inactive You can't update the show")
	}

	result := databaseauth.DB.Model(&existingShow).Updates(updatedFields)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
