package databaseauth

import (
	"booking-service/models"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
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

	var config Config

	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
	}
	// Construct the path to the config file
	configPath := filepath.Join(dir, "config.yaml")

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
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

	err = DB.AutoMigrate(&models.Booking{}) // Automigrate the schema
	if err != nil {
		panic("Automigration failed")
	}

}
