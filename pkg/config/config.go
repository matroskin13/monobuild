package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"monobuild/internal/errors"
	"monobuild/internal/files"
	"path"
)

const DefaultConfigName = ".monobuild.yml"

type Config struct {
	Packages map[string]Package
}

type Package struct {
	Entry     string
	FullEntry string
	Build     struct {
		Docker *struct {
			Image string
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

func ParseDefaultConfig(applicationPath string) (*Config, error) {
	configPath := path.Join(applicationPath, DefaultConfigName)

	if !files.FileExists(configPath) {
		return nil, errors.NewRichError(fmt.Sprintf("Configuration file not found in %s, please specity correct path", configPath), nil)
	}

	return ParseConfigFromFile(path.Join(applicationPath, DefaultConfigName))
}
