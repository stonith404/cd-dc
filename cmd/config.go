package main

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
	Services map[string]Service `yaml:"services"`
	ApiKey   string             `yaml:"apiKey"`
}

func Get() Config {
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
	return Get().ApiKey
}

func GetServices() map[string]Service {
	return Get().Services
}
