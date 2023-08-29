package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Path       string   `yaml:"path"`
	Containers []string `yaml:"containers"`
}

type Config struct {
	Services             map[string]Service `yaml:"services"`
	ApiKey               string             `yaml:"apiKey"`
	NumberOfImagesToKeep int                `yaml:"numberOfImagesToKeep"`
}

func get() Config {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// create a config struct and deserialize the data into that struct
	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return config
}

func GetApiKey() string {
	return get().ApiKey
}

func GetServices() map[string]Service {
	return get().Services
}

func GetNumberOfImagesToKeep() int {
	nuberOfImagesToKeep := get().NumberOfImagesToKeep
	if nuberOfImagesToKeep == 0 {
		return 1
	}
	return nuberOfImagesToKeep
}
