package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

func NewConfig() (config *Config) {
	config = &Config{}

	flag.StringVar(&config.ServerAddress, "a", ":8080", "port to listen on")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&config.FileStoragePath, "f", "", "path to file where to store data")

	serverAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		config.ServerAddress = serverAddress
	}
	baseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		config.BaseURL = baseURL
	}
	fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		config.FileStoragePath = fileStoragePath
	}

	flag.Parse()

	return config
}
