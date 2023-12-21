package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	BookingID  string         `gorm:"column:booking_id;primaryKey" json:"bookingId"`
	UserID     int64          `gorm:"type:BIGINT(7);column:user_id;not null" json:"userId"`
	MovieID    int64          `gorm:"type:BIGINT(7);column:movie_id;not null" json:"movieId"`
	TheatreID  int64          `gorm:"type:BIGINT(7);column:theatre_id;not null" json:"theatreId"`
	ShowID     int64          `gorm:"type:BIGINT(7);column:show_id;not null" json:"showId"`
	ShowTime   datatypes.Time `gorm:"column:show_time;not null" json:"showTime"`
	ShowDate   datatypes.Date `gorm:"column:show_date;not null" json:"showDate"`
	Tickets    int            `gorm:"column:tickets;not null" json:"tickets"`
	BookedTime time.Time      `gorm:"column:booked_time;not null" json:"bookedTime"`
	Status     BookingStatus  `gorm:"column:status;type:ENUM('pending','confirmed','canceled');"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}
