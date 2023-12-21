package databaseauth

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"user-service/models"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Database struct {
		Host     string
		Port     int
		UserName string
		Password string
		DBName   string
	}
}

func InitDB() {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Println("Error in getting Absolute Path")
	}
	configPath := filepath.Join(dir, "config.yaml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error Reading configuration files: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Error decoding configuration")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&&parseTime=True&loc=Local",
		config.Database.UserName,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect to database")
	}
	DB = db

	err = DB.AutoMigrate(&models.User{}, &models.PasswordReset{})
	if err != nil {
		log.Panic("Automigration failed")
	}

}
