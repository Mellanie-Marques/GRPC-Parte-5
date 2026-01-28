package config

import "os"

func GetApplicationPort() string {
	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "3002"
	}
	return port
}
