package config

import (
	"encoding/json"
	"os"

	"golang-mock/model"
)

// LoadConfigs ...
func LoadConfigs(path string) ([]model.MockConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var configs []model.MockConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

// SaveConfigs ...
func SaveConfigs(path string, configs []model.MockConfig) error {
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
