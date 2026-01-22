package config

import (
	"os"
	"path/filepath"
	"testing"
)

const validConfigYAML = `
aws:
  profile: hacc-user
  kms_id: hacc-key
  param_path: hacc-vault

security:
  create_scp: false
  scp_name: ""
  member_role: ""

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
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := Load("/path/does/not/exist.yaml")
	if err == nil {
		t.Fatalf("expected error for missing config file")
	}
}

func TestSecurityValidation_Disabled(t *testing.T) {
	path := writeTempConfig(t, validConfigYAML)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestSecurityValidation_EnabledMissingFields(t *testing.T) {
	invalidYAML := `
security:
  create_scp: true
  scp_name: ""
  member_role: ""
`

	path := writeTempConfig(t, invalidYAML)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error when create_scp=true")
	}
}
