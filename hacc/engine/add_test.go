package engine

import (
	"os"
	"testing"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/config"
	"github.com/nbailey20/hacc/vault"
	"github.com/stretchr/testify/require"
)

func TestAddCredential(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/add",
		ObfuscationKey: "secretkey",
	}
	cmd := cli.CLICommand{
		Action:   cli.AddAction{},
		Service:  "testservice",
		Username: "testuser",
		Password: "testpass",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Test adding a credential
	result := testExecutor.Execute(cmd)
	require.True(t, result.Success)
	require.Equal(t, "add", result.Action)
	require.Len(t, result.Data, 1)
	require.Equal(t, "testservice", result.Data[0].Service)
	require.Equal(t, "testuser", result.Data[0].Username)

	// Clean up
	testVault.Delete("testservice", "testuser")
}

func TestAddMultiCredential(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/multi",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary Executor with a temporary vault
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Create a temporary JSON file with multiple credentials
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	testData := `{"creds_list":[{"service":"github","username":"user1","password":"pass1"},{"service":"gitlab","username":"user2","password":"pass2"}]}`
	_, err = tmpFile.WriteString(testData)
	require.NoError(t, err)
	tmpFile.Close()

	// Test adding multiple credentials using a file
	cmd := cli.CLICommand{
		Action: cli.AddAction{},
		File:   tmpFile.Name(),
	}
	result := testExecutor.Execute(cmd)

	// Verify the result indicates success
	require.True(t, result.Success)
	require.Equal(t, "add", result.Action)
	require.Len(t, result.Data, 2)

	// Verify expected credentials are in the result data
	expectedCreds := map[string]string{
		"github": "user1",
		"gitlab": "user2",
	}
	for _, cred := range result.Data {
		expectedUser, ok := expectedCreds[cred.Service]
		require.True(t, ok, "unexpected credential in result: %s", cred.Service)
		require.Equal(t, expectedUser, cred.Username)
	}

	// Clean up
	testVault.Delete("github", "user1")
	testVault.Delete("gitlab", "user2")
}

func TestAddMultiCredentialFileNotFound(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/multi",
		ObfuscationKey: "secretkey",
	}
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Test case: File doesn't exist
	result := testExecutor.addMultiCredential("nonexistent.json")
	require.False(t, result.Success)
	require.Equal(t, "add", result.Action)
	require.Contains(t, result.Error, "open nonexistent.json: The system cannot find the file specified.")
}

func TestAddMultiCredentialInvalidJson(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/multi",
		ObfuscationKey: "secretkey",
	}
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Create a temporary file with invalid JSON
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("invalid json {")
	tmpFile.Close()

	// Test case: Invalid JSON
	result := testExecutor.addMultiCredential(tmpFile.Name())
	require.False(t, result.Success)
	require.Equal(t, "add", result.Action)
	require.Contains(t, result.Error, "invalid character")
}
