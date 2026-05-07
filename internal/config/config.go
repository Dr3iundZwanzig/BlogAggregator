package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file")
	}
	var configFile Config
	err = json.Unmarshal(file, &configFile)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshal file")
	}
	return configFile, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting config file path")
	}
	return home + "/" + configFilename, nil
}
func write(configStruct Config) error {
	data, err := json.Marshal(configStruct)
	if err != nil {
		return fmt.Errorf("error writing file")
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
