package utils

import (
	"errors"

	ini "gopkg.in/ini.v1"
)

var global_config *Config

func GetConfig() (*Config, error) {
	if global_config != nil {
		return nil, errors.New("not ready yet")
	}
	return global_config, nil
}

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
	global_config = config
	return config, nil
}
