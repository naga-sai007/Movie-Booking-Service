package databaseauth

import (
	"fmt"
	"log"
	"movie-service/models"
	"os"
	"path/filepath"

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

	err = DB.AutoMigrate(&models.Movie{}) // Automigrate the schema
	if err != nil {
		panic("Automigration failed")
	}

}
