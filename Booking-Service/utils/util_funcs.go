package utils

import (
	"booking-service/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

type Config struct {
	Showservice struct {
		Geturl    string
		Updateurl string
	}
	Userservice struct {
		Userurl string
	}
	SecretKey string `yaml:"secretKey"`
	ApiKey    string `yaml:"apiKey"`
}

const (
	bookingIDPrefix = "BO"
	letters         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits          = "0123456789"
)

type ShowData struct {
	ShowId         int64
	MovieId        int64
	TheatreId      int64
	ShowTime       datatypes.Time
	ShowDate       datatypes.Date
	ShowStatus     bool
	AvailableSeats int
	TicketsSold    int
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func ValidateBooking(NewBooking models.Booking) (ShowData, error) {
	showId := NewBooking.ShowID
	var showData ShowData
	showData, err := GetShowData(showId)
	if err != nil {
		fmt.Println(err)
		return ShowData{}, ValidationError{Message: "failed to get show data"}
	}
	if !showData.ShowStatus {
		return ShowData{}, ValidationError{Message: "Sorry Show Started or Show is Inactive,Booking for this show stopped"}
	}

	if NewBooking.Tickets > showData.AvailableSeats && showData.AvailableSeats == 0 {
		return ShowData{}, ValidationError{Message: "Not enough available seats"}
	}

	if NewBooking.MovieID != showData.MovieId || NewBooking.TheatreID != showData.TheatreId {
		return ShowData{}, ValidationError{Message: "Movie-Id or theatre-ID mismatch with show-ID"}
	}

	if NewBooking.Tickets <= 0 || NewBooking.Tickets > showData.AvailableSeats {
		return ShowData{}, ValidationError{Message: "Invalid number of tickets"}
	}
	if !showData.ShowStatus {
		return ShowData{}, ValidationError{Message: "Sorry Booking for this show Stopped"}
	}

	return showData, nil
}

func GetShowData(showId int64) (ShowData, error) {
	config, err := GetConfigFromConfigFile()
	if err != nil {
		return ShowData{}, err
	}
	baseUrl, err := GetURLFromConfig(config, "base")
	if err != nil {
		return ShowData{}, err
	}
	url := fmt.Sprintf("%s%d", baseUrl, showId)
	response, err := http.Get(url)
	if err != nil {
		return ShowData{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ShowData{}, fmt.Errorf("HTTP request failed with status: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ShowData{}, err
	}

	var showData ShowData
	err = json.Unmarshal(body, &showData)
	if err != nil {
		return ShowData{}, err
	}
	return showData, nil
}

func GetConfigFromConfigFile() (Config, error) {
	var config Config

	viper.SetConfigFile("C:/Users/NAGASAIA/GoProjects/movie-booking-system/booking-service/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %s", err)
	}

	return config, nil
}

func GetURLFromConfig(config Config, key string) (string, error) {
	switch key {
	case "base":
		return config.Showservice.Geturl, nil
	case "update":
		return config.Showservice.Updateurl, nil
	case "user":
		return config.Userservice.Userurl, nil
	default:
		return "", fmt.Errorf("unknown key: %s", key)
	}
}

func UpdateTicketsInShowTable(showID int64, numberOfTickets int, existingSeats int, ticketsSold int) error {

	config, err := GetConfigFromConfigFile()
	if err != nil {
		fmt.Println(err)
		return err
	}
	updateUrl, err := GetURLFromConfig(config, "update")
	if err != nil {
		fmt.Println(err)
		return err
	}
	url := fmt.Sprintf("%s%d", updateUrl, showID)

	updated_available_seats := existingSeats - numberOfTickets
	if existingSeats < numberOfTickets && updated_available_seats < 0 {
		return err
	}
	updated_tickets_sold := ticketsSold + numberOfTickets
	//fmt.Println(updated_tickets_sold)

	updateRequest := map[string]interface{}{
		"show_id":         showID,
		"available_seats": updated_available_seats,
		"tickets_sold":    updated_tickets_sold,
	}

	requestBody, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	ApiKey := config.ApiKey
	AuthToken, err := HashSecurityKey(ApiKey)
	if err != nil {
		return err
	}
	ApikeyWithPrefix := "Apikey " + AuthToken
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ApikeyWithPrefix)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update show. Status code: %d", resp.StatusCode)
	}
	return nil
}

func ValidateUserId(userId int64, authToken string) error {

	config, err := GetConfigFromConfigFile()
	if err != nil {
		return err
	}
	userUrl, err := GetURLFromConfig(config, "user")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%d", userUrl, userId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set the Authorization header with the authentication token
	req.Header.Set("Authorization", authToken)

	// Use http.DefaultClient to send the request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status: %s", response.Status)
	}

	return nil
}

func GenerateBookingID() string {
	rand.Seed(time.Now().UnixNano())

	// Generate random letters
	var randomLetters string
	for i := 0; i < 2; {
		randomLetters += string(letters[rand.Intn(len(letters))])
		i++
	}

	// Generate random digits
	randomDigits := fmt.Sprintf("%03d", rand.Intn(1000))

	// Concatenate the prefix, letters, and digits
	return fmt.Sprintf("%s%s%s", bookingIDPrefix, randomLetters, randomDigits)
}

func HashSecurityKey(ApiKey string) (string, error) {
	// Hash the user's password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ApiKey), bcrypt.DefaultCost)
	if err != nil {
		return " ", err
	}
	return string(hashedPassword), nil
}
