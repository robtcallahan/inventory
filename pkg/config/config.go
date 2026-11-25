package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type DBConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Config ...
type Config struct {
	Database     DBConfig `json:"database"`
	TestDatabase DBConfig `json:"test_database"`
	PhotosDir string `json:"photos_dir"`
	WebDir string `json:"web_dir"`
}

// ReadConfig ...
func ReadConfig(file string) (*Config, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config = Config{}
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config file: %w", err)
	}
	if config.PhotosDir == "" { config.PhotosDir = "photos" }
	if config.WebDir == "" { config.WebDir = "./web/public" }
	return &config, nil
}
