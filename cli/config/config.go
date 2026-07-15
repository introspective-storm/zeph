package config

import (
	"errors"
	"os"
	"time"

	"go.yaml.in/yaml/v4"
)

type ProjectConfig struct {
	Name       string    `yaml:"name"`
	CreatedAt  time.Time `yaml:"created_at"`
	LastOpened time.Time `yaml:"last_opened"`
	DataSource string    `yaml:"data_source"`
	ModelPath  string    `yaml:"model_path"`
	Tests      []string  `yaml:"tests"` // Update this line!
}

func (c *ProjectConfig) Save(filePath string) error {
	c.LastOpened = time.Now()

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func Load(filePath string) (*ProjectConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var c ProjectConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	if c.Name == "" {
		return nil, errors.New("missing project name: required metadata")
	}

	return &c, nil
}
