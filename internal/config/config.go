package config

import (
	"os"
	"fmt"
	"encoding/json"
	"path/filepath"
)

type Config struct {
	URL			string	`json:"db_url"`
	CurrentUser	string	`json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error retrieving user home directory: ", err)
	}

	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {
	var cfg Config

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	jsonFile, err := os.Open(configFilePath)
    if err != nil {
        return cfg, fmt.Errorf("Error opening file:", err)
    }
    defer jsonFile.Close()

    byteValue, err := os.ReadFile(configFilePath)
    if err != nil {
        return cfg, fmt.Errorf("Error reading file:", err)
    }

    err = json.Unmarshal(byteValue, &cfg)
    if err != nil {
        return cfg, fmt.Errorf("Error unmarshaling JSON:", err)
    }

    return cfg, nil
}

func write(cfg Config) error {
	// Marshal the struct into a pretty-printed JSON byte slice
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling to JSON:", err)
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// os.WriteFile creates the file if it doesn't exist,
	// or truncates it if it does, then writes the content.
	err = os.WriteFile(configFilePath, jsonData, 0644) // 0644 sets file permissions
	if err != nil {
		return fmt.Errorf("Error writing to file: ", err)
	}

	return nil
}

func (c Config) SetUser(user string) error {
	c.CurrentUser = user
	err := write(c)
	if err != nil {
		return err
	}

	return nil
}