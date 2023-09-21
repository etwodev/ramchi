package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const CONFIG = "./ramchi.config.json"

var c *Config

func Load() error {
	_, err := os.Stat(CONFIG)
	if os.IsNotExist(err) {
		if err := Create(); err != nil {
			return fmt.Errorf("Load: failed creating load: %w", err)
		}
	}

	file, err := os.ReadFile(CONFIG)
	if err != nil {
		return fmt.Errorf("Load: failed reading json: %w", err)
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		return fmt.Errorf("Load: failed marshalling json: %w", err)
	}
	return nil
}

func Create() error {
	file, err := json.MarshalIndent(&Config{Port: "8080", Address: "localhost", Experimental: false}, "", " ")
	if err != nil {
		return fmt.Errorf("Create: failed marshalling config: %w", err)
	}
	err = os.WriteFile(CONFIG, file, 0644)
	if err != nil {
		return fmt.Errorf("Create: failed writing config: %w", err)
	}
	return nil
}

func New() error {
	if c == nil {
		err := Load()
		if err != nil {
			return fmt.Errorf("New: failed loading json: %w", err)
		}
	}
	return nil
}
