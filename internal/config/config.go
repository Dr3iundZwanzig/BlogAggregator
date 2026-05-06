package config

import (
	"encoding/json"
	"os"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(name string) {
	c.CurrentUserName = name
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}
	var configFile Config
	err = json.Unmarshal(file, &configFile)
	if err != nil {
		return Config{}, err
	}
	return configFile, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + configFilename, nil
}
