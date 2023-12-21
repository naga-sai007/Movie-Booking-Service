package models

import (
	"gorm.io/datatypes"
)

type Shows struct {
	ShowId         int64          `gorm:"type:BIGINT(7);column:show_id;primaryKey" json:"showId"`
	MovieID        int64          `gorm:"type:BIGINT(7);column:movie_id" json:"movieId"`
	TheatreID      int64          `gorm:"type:BIGINT(7);column:theatre_id;foreignKey" json:"theatreId"`
	ShowTime       datatypes.Time `gorm:"column:show_time;type:time" json:"showTime"`
	ShowDate       datatypes.Date `gorm:"column:show_date;type:date" json:"showDate"`
	ShowStatus     bool           `gorm:"column:show_status;default:1" json:"showStatus"`
	AvailableSeats int            `gorm:"column:available_seats" json:"availableSeats"`
	TicketsSold    int            `gorm:"column:tickets_sold" json:"ticketsSold"`
}
