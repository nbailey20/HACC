package display

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nbailey20/hacc/hacc/config"
	vault "github.com/nbailey20/hacc/hacc/vault"
)

func TestAddMultiCredential(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath: "/hackyclient/test/multi/",
	}
	// Create a temporary vault
	testVault, err := vault.NewVault(nil, cfg)
	require.NoError(t, err)

	// Create a temporary JSON file with multiple credentials
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	testData := `{"creds_list":[{"service":"github","username":"user1","password":"pass1"},{"service":"gitlab","username":"user2","password":"pass2"}]}`
	_, err = tmpFile.WriteString(testData)
	require.NoError(t, err)
	tmpFile.Close()

	// Test case 1: Successfully add multiple credentials
	addMultiCmd := addMultiCredentialCmd(*testVault, tmpFile.Name())
	msg := addMultiCmd()

	// Should return AddSuccessMsg with no error if all succeed
	_, ok := msg.(AddSuccessMsg)
	require.True(t, ok)

	// Verify credentials were added
	cred1, err := testVault.Get("github", "user1")
	require.NoError(t, err)
	require.Equal(t, "pass1", cred1)

	cred2, err := testVault.Get("gitlab", "user2")
	require.NoError(t, err)
	require.Equal(t, "pass2", cred2)

	// Clean up
	testVault.Delete("github", "user1")
	testVault.Delete("gitlab", "user2")
}

func TestAddMultiCredentialFileNotFound(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath: "/hackyclient/test/multi/",
	}
	testVault, err := vault.NewVault(nil, cfg)
	require.NoError(t, err)

	// Test case: File doesn't exist
	addMultiCmd := addMultiCredentialCmd(*testVault, "nonexistent.json")
	msg := addMultiCmd()

	failMsg, ok := msg.(AddFailedMsg)
	require.True(t, ok)
	require.Error(t, failMsg.Error)
}

func TestAddMultiCredentialInvalidJson(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath: "/hackyclient/test/multi/",
	}
	testVault, err := vault.NewVault(nil, cfg)
	require.NoError(t, err)

	// Create a temporary file with invalid JSON
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("invalid json {")
	tmpFile.Close()

	// Test case: Invalid JSON
	addMultiCmd := addMultiCredentialCmd(*testVault, tmpFile.Name())
	msg := addMultiCmd()

	failMsg, ok := msg.(AddFailedMsg)
	require.True(t, ok)
	require.Error(t, failMsg.Error)
}
