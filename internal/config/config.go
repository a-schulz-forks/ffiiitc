package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	FireflyAppTimeout = 10               // 10 sec for fftc to app service timeout
	ModelFile         = "data/model.gob" //file name to store model
)

type Config struct {
	APIKey         string
	FFApp          string
	ClassifierPort int
}

var envVars = []string{
	"FF_API_KEY",
	"FF_APP_URL",
}

func EnvVarExist(varName string) bool {
	_, present := os.LookupEnv(varName)
	return present
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewConfig() (*Config, error) {
	for _, val := range envVars {
		exist := EnvVarExist(val)
		if !exist {
			return nil, errors.New("env variable is not set: " + val)
		}
	}
	port, err := strconv.Atoi(getEnv("APPLICATION_PORT", "8080"))
	if err != nil {
		// Handle the error, e.g., log it or return a default value
		fmt.Printf("Error converting port to integer: %s\n", err)
	}

	cfg := Config{
		APIKey:         os.Getenv("FF_API_KEY"),
		FFApp:          os.Getenv("FF_APP_URL"),
		ClassifierPort: port,
	}

	return &cfg, nil
}
