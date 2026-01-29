package config

import (
	"os"
	"path/filepath"
	"testing"
)

const validConfigYAML = `
aws:
  profile: hacc-user
  kms_id: 1234567890
  param_path: /hacc-vault
  obfuscation_key: supersecretkey

client:
  check_for_upgrades: true
  cleanup_old_versions: false
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

	if !cfg.Client.CheckForUpgrades {
		t.Errorf("expected check_for_upgrades=true")
	}
	if cfg.Client.CleanupOldVersions {
		t.Errorf("expected cleanup_old_versions=false")
	}
	if cfg.AWS.Profile != "hacc-user" {
		t.Errorf("expected aws.profile=hacc-user, got %s", cfg.AWS.Profile)
	}
	if cfg.AWS.KmsId != "1234567890" {
		t.Errorf("expected aws.kms_id=hacc-key, got %s", cfg.AWS.KmsId)
	}
	if cfg.AWS.ParamPath != "/hacc-vault" {
		t.Errorf("expected aws.param_path=hacc-vault, got %s", cfg.AWS.ParamPath)
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
aws:
  kms_id: 1234567890
  param_path: /hacc-vault
  obfuscation_key: supersecretkey
client:
  check_for_upgrades: true
  cleanup_old_versions: false
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
aws:
  profile: hacc-user
  param_path: /hacc-vault
  obfuscation_key: supersecretkey
client:
  check_for_upgrades: true
  cleanup_old_versions: false
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
aws:
  profile: hacc-user
  kms_id: 1234567890
  obfuscation_key: supersecretkey	
client:
  check_for_upgrades: true
  cleanup_old_versions: false
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
aws:
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
