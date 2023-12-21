package utils

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	SecretKey string `yaml:"secretKey"`
}

func GetSecretKeyFromConfig() string {
	cfg, err := ReadConfig()
	if err != nil {
		// Handle the error (e.g., log, return an error, etc.)
		fmt.Println("Error reafing config")
		return ""
	}
	// Access the secret key
	secretKey := cfg.SecretKey
	return secretKey
}

func ReadConfig() (Config, error) {
	var cfg Config
	// // Set the file name and path
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
	}
	// Construct the path to the config file
	//configPath := filepath.Join(dir, "config.yaml")
	//configPath = filepath.ToSlash(configPath)
	viper.AddConfigPath(dir) // You may adjust the path as needed
	//fmt.Println("Raw Config Path:", configPath)

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return cfg, err
	}

	// Unmarshal the configuration into the struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
