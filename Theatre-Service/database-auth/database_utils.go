package databaseauth

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"theatre-service/models"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		DBName   string
	}
}

func InitDb() {

	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
	}
	// Construct the path to the config file
	configPath := filepath.Join(dir, "config.yaml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	// Unmarshal the configuration values into a Config struct
	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Error decoding configuration: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName)

	//dsn := "root:12345678@tcp(localhost:3306)/movie-booking-service?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	DB = db

	err = DB.AutoMigrate(&models.Theatre{}) // Automigrate the schema
	if err != nil {
		panic("Automigration failed")
	}

	err = DB.AutoMigrate(&models.Shows{}) // Automigrate the schema
	if err != nil {
		panic("Automigration failed")
	}

	eventName := "update_show_status"
	if !EventExists(DB, eventName) {
		// Create the event if it doesn't exist
		DB.Exec("SET GLOBAL event_scheduler = ON")
		createEventQuery := `
			CREATE EVENT update_show_status
			ON SCHEDULE EVERY 15 MINUTE
			DO
				UPDATE shows
				SET active_status = 0
				WHERE show_date < CURRENT_DATE
					OR (show_date = CURRENT_DATE AND show_time < CURRENT_TIME);
		`
		DB.Exec(createEventQuery)
		fmt.Println("Event created successfully.")
	} else {
		fmt.Println("Event already exists.")
	}

}

func EventExists(db *gorm.DB, eventName string) bool {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.EVENTS WHERE EVENT_NAME = ?", eventName).Scan(&count)
	return count > 0
}
