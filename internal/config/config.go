package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("error: %v", err)
	}

	file_path := filepath.Join(homeDir, configFileName)

	data, err := os.ReadFile(file_path)
	if err != nil {
		return Config{}, fmt.Errorf("error: %v", err)
	}

	var json_file_struct Config

	err = json.Unmarshal(data, &json_file_struct)
	if err != nil {
		return json_file_struct, fmt.Errorf("error: %v", err)
	}

	return json_file_struct, nil
}

func (c *Config) SetUser(user_name string) error {
	c.CurrentUserName = user_name

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	file_path := filepath.Join(homeDir, configFileName)

	err = os.WriteFile(file_path, data, 0644)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	return nil
}
