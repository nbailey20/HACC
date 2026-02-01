package config

import (
	"os"
	"path/filepath"
	"testing"
)

const validConfigYAML = `
profile: hacc-user
kms_id: 1234567890
param_path: /hacc-vault
obfuscation_key: supersecretkey
`

func writeTempConfig(t *testing.T, contents string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	return path
}

func TestLoadConfig_Success(t *testing.T) {
	path := writeTempConfig(t, validConfigYAML)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Profile != "hacc-user" {
		t.Errorf("expected profile=hacc-user, got %s", cfg.Profile)
	}
	if cfg.KmsId != "1234567890" {
		t.Errorf("expected kms_id=hacc-key, got %s", cfg.KmsId)
	}
	if cfg.ParamPath != "/hacc-vault" {
		t.Errorf("expected param_path=hacc-vault, got %s", cfg.ParamPath)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := Load("/path/does/not/exist.yaml")
	if err == nil {
		t.Fatalf("expected error for missing config file")
	}
}

func TestValidate(t *testing.T) {
	path := writeTempConfig(t, validConfigYAML)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidate_MissingAWSProfile(t *testing.T) {
	path := writeTempConfig(t, `
kms_id: 1234567890
param_path: /hacc-vault
obfuscation_key: supersecretkey
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected validation success for missing aws.profile")
	}
}

func TestValidate_MissingKmsId(t *testing.T) {
	path := writeTempConfig(t, `
profile: hacc-user
param_path: /hacc-vault
obfuscation_key: supersecretkey
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected validation success with no aws.kms_id")
	}
}

func TestValidate_MissingParamPath(t *testing.T) {
	path := writeTempConfig(t, `
profile: hacc-user
kms_id: 1234567890
obfuscation_key: supersecretkey	
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error for missing aws.param_path")
	}
}

func TestValidate_MissingObfuscationKey(t *testing.T) {
	path := writeTempConfig(t, `
profile: hacc-user
kms_id: 1234567890
param_path: /hacc-vault
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error for missing aws.obfuscation_key")
	}
}
