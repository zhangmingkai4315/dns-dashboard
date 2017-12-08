package utils

import (
	"errors"
	"os"

	ini "gopkg.in/ini.v1"
)

var globalConfig *Config

func init() {
	// if config exist then load the config.ini
	if _, err := os.Stat("../config.ini"); err == nil {
		globalConfig, _ = LoadConfigFromFile("../config.ini")
	}
}

// GetConfig will get the global config from module
func GetConfig() (*Config, error) {
	if globalConfig != nil {
		return nil, errors.New("config not ready yet")
	}
	return globalConfig, nil
}

// LoadConfigFromFile will load config from file
func LoadConfigFromFile(file string) (*Config, error) {
	cfg, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = cfg.MapTo(config)
	if err != nil {
		return nil, err
	}
	globalConfig = config
	return config, nil
}
