package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AWS      AWSConfig      `yaml:"aws"`
	Security SecurityConfig `yaml:"security"`
	Client   ClientConfig   `yaml:"client"`
}

type AWSConfig struct {
	Profile   string `yaml:"profile"`
	KmsId     string `yaml:"kms_id"`
	ParamPath string `yaml:"param_path"`
}

type SecurityConfig struct {
	CreateSCP  bool   `yaml:"create_scp"`
	SCPName    string `yaml:"scp_name"`
	MemberRole string `yaml:"member_role"`
}

type ClientConfig struct {
	CheckForUpgrades   bool `yaml:"check_for_upgrades"`
	CleanupOldVersions bool `yaml:"cleanup_old_versions"`
}

func GetPath() (string, error) {
	// returns path of config.yaml file required for program
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hacc", "config.yaml"), nil
}

func (c Config) Validate() error {
	requiredStrings := map[string]string{
		"aws.profile":    c.AWS.Profile,
		"aws.param_path": c.AWS.ParamPath,
	}
	for name, value := range requiredStrings {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required in config.yaml.", name)
		}
	}
	return nil
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
