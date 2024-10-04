package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

type ConfigLoader struct {
	Config any
}

func NewConfigLoader(config any) *ConfigLoader {
	return &ConfigLoader{Config: config}
}

func (cl *ConfigLoader) LoadConfig(filePath string) error {
	switch ext := getFileExtension(filePath); ext {
	case ".yaml", ".yml":
		return cl.LoadYAML(filePath)
	case ".json":
		return cl.LoadJSON(filePath)
	case ".ini":
		return cl.LoadINI(filePath)
	case ".toml":
		return cl.LoadTOML(filePath)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

func (cl *ConfigLoader) LoadYAML(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}
	err = yaml.Unmarshal(data, &cl.Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return nil
}

func (cl *ConfigLoader) LoadJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}
	err = json.Unmarshal(data, &cl.Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

func (cl *ConfigLoader) LoadINI(filePath string) error {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return fmt.Errorf("failed to read INI file: %w", err)
	}
	err = cfg.MapTo(&cl.Config)
	if err != nil {
		return fmt.Errorf("failed to map INI to config: %w", err)
	}
	return nil
}

func (cl *ConfigLoader) LoadTOML(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read TOML file: %w", err)
	}
	err = toml.Unmarshal(data, &cl.Config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal TOML: %w", err)
	}
	return nil
}

func getFileExtension(filePath string) string {
	return filepath.Ext(filePath)
}
