package repository

import (
	databaseauth "theatre-service/database-auth"
	"theatre-service/models"
)

func CreateTheatre(newTheatre models.Theatre) {

	databaseauth.DB.Create(&newTheatre)
}

func GetAllTheatres() []models.Theatre {
	var theatres []models.Theatre
	databaseauth.DB.Preload("Shows").Find(&theatres)
	return theatres
}

func GetTheatreByID(theatreID string) (*models.Theatre, error) {
	var Theatre models.Theatre
	if err := databaseauth.DB.Preload("Shows").Find(&Theatre, theatreID).Error; err != nil {
		return nil, err
	}
	return &Theatre, nil
}
