package repository

import (
	databaseauth "booking-service/database-auth"
	"booking-service/models"
	"booking-service/utils"
	"fmt"
	"log"
	"time"
)

func CreateBooking(NewBooking models.Booking, authToken string) (models.Booking, error) {
	showData, validationErr := utils.ValidateBooking(NewBooking)
	if validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return models.Booking{}, validationErr
	}

	userId := NewBooking.UserID
	if err := utils.ValidateUserId(userId, authToken); err != nil {
		return models.Booking{}, err
	}

	showId := NewBooking.ShowID
	numberOfTickets := NewBooking.Tickets
	//fmt.Println(showData.TicketsSold)
	if err := utils.UpdateTicketsInShowTable(showId, numberOfTickets, showData.AvailableSeats, showData.TicketsSold); err != nil {
		return models.Booking{}, err
	}

	NewBooking.Status = "confirmed"
	NewBooking.BookedTime = time.Now()
	NewBooking.ShowTime = showData.ShowTime
	NewBooking.ShowDate = showData.ShowDate
	NewBooking.BookingID = utils.GenerateBookingID()

	if err := databaseauth.DB.Create(&NewBooking).Error; err != nil {
		log.Printf("Error creating booking: %v", err)
		return models.Booking{}, err
	}

	return NewBooking, nil
}

func GetBookingById(BookingId string) (*models.Booking, error) {
	var Booking models.Booking
	if err := databaseauth.DB.Where("booking_id = ?", BookingId).First(&Booking).Error; err != nil {
		return &models.Booking{}, err
	}
	return &Booking, nil
}

func GetBookingsByUserId(userId int64, authToken string) ([]models.Booking, error) {

	err := utils.ValidateUserId(userId, authToken)
	if err != nil {
		fmt.Println(err)
		return []models.Booking{}, err
	}
	var bookings []models.Booking
	result := databaseauth.DB.Where("user_id=?", userId).Find(&bookings)
	if result.Error != nil {
		fmt.Println(err)
		return []models.Booking{}, err
	}
	return bookings, nil
}

// func UpdateBooking(updatedBooking models.Booking) {
// 	bookingId := updatedBooking.BookingID
// 	var existingBooking models.Booking
// 	databaseauth.DB.First(&existingBooking, bookingId)
// 	existingBooking.MovieID = updatedBooking.MovieID
// 	existingBooking.TheatreID = updatedBooking.TheatreID
// 	existingBooking.ShowTime = updatedBooking.ShowTime
// 	existingBooking.Status = updatedBooking.Status
// 	databaseauth.DB.Save(&existingBooking)
// }

// func CancelBooking() {}

// func DeleteBooking(bookingID string) {
// 	databaseauth.DB.Delete(&models.Booking{}, bookingID)
// }
