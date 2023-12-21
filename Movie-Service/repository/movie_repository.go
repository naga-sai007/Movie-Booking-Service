package repository

import (
	databaseauth "movie-service/database-auth"
	"movie-service/models"
)

func CreateMovie(newMovie models.Movie) (models.Movie, error) {

	result := databaseauth.DB.Create(&newMovie)
	if result.Error != nil {
		return models.Movie{}, result.Error
	}
	return newMovie, nil
}

func GetAllMovies() []models.Movie {
	var movies []models.Movie
	databaseauth.DB.Find(&movies)
	return movies
}

func GetMovieByID(MovieID string) (*models.Movie, error) {
	var Movie models.Movie
	if err := databaseauth.DB.First(&Movie, MovieID).Error; err != nil {
		return nil, err
	}
	return &Movie, nil
}

func DeleteMovie(MovieID string) {

	databaseauth.DB.Delete(&models.Movie{}, MovieID)
}
