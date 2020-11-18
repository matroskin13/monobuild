package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Packages map[string]struct {
		Entry string
		Build struct {
			Docker struct {
				Image string
			}
		}
	}
}

func ParseConfigFromFile(path string) (*Config, error) {
	var cfg Config

	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file to Config: %w", err)
	}

	return &cfg, nil
}
