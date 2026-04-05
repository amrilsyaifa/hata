package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFile = "i18n.config.yml"

type Config struct {
	ProjectID string    `yaml:"project_id"`
	Sheet     SheetConf `yaml:"sheet"`
	Auth      AuthConf  `yaml:"auth"`
	Languages []string  `yaml:"languages"`
	Paths     PathConf  `yaml:"paths"`
	Options   Options   `yaml:"options"`
}

type SheetConf struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type AuthConf struct {
	Type            string `yaml:"type"`
	CredentialsPath string `yaml:"credentials_path"`
	TokenPath       string `yaml:"token_path"`
}

type PathConf struct {
	Base   string `yaml:"base"`
	Output string `yaml:"output"`
}

type Options struct {
	NestedJSON bool `yaml:"nested_json"`
	SortKeys   bool `yaml:"sort_keys"`
	KeepUnused bool `yaml:"keep_unused"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
