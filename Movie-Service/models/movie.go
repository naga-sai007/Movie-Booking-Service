package models

type Movie struct {
	MovieID     int64   `gorm:"type:BIGINT(7);column:movie_id;primaryKey" json:"movieId"`
	Title       string  `gorm:"column:movie_name" json:"title"`
	Genre       string  `gorm:"column:genre" json:"genre"`
	Director    string  `gorm:"column:director" json:"director"`
	UserRatings float64 `gorm:"column:user_ratings" json:"userRatings"`
	Cast        string  `gorm:"column:cast" json:"cast"`
	Plot        string  `gorm:"column:plot" json:"plot"`
}
