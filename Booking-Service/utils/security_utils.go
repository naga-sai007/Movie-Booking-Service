package utils

import (
	"fmt"
)

func GetSecretKeyFromConfig() string {
	cfg, err := GetConfigFromConfigFile()
	if err != nil {
		// Handle the error (e.g., log, return an error, etc.)
		fmt.Println("Error reafing config")
		return ""
	}
	// Access the secret key
	secretKey := cfg.SecretKey
	return secretKey
}
