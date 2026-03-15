package engine

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/config"
	"github.com/nbailey20/hacc/helpers"
	"github.com/nbailey20/hacc/vault"
	"github.com/stretchr/testify/require"
)

func TestBackupCredential(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/backup",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Add a test credential
	err = testVault.Add("backupservice", "backupuser", "backuppass")
	require.NoError(t, err)

	// Test backing up the credential to a file
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	cmd := cli.CLICommand{
		Action:   cli.BackupAction{},
		File:     tmpFile.Name(),
		Service:  "backupservice",
		Username: "backupuser",
	}
	result := testExecutor.Execute(cmd)
	require.True(t, result.Success)
	require.Equal(t, "backup", result.Action)
	require.Len(t, result.Data, 1)
	require.Equal(t, "backupservice", result.Data[0].Service)
	require.Equal(t, "backupuser", result.Data[0].Username)

	// Clean up
	testVault.Delete("backupservice", "backupuser")
}

func TestBackupServiceCredentials(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/backup",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Add test credentials
	err = testVault.Add("svc1", "user1", "pass1")
	require.NoError(t, err)
	err = testVault.Add("svc1", "user2", "pass2")
	require.NoError(t, err)

	// Test backing up all credentials for a service
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	cmd := cli.CLICommand{
		File:    tmpFile.Name(),
		Service: "svc1",
	}
	result := testExecutor.backup(cmd)
	require.True(t, result.Success)
	require.Equal(t, "backup", result.Action)
	require.Len(t, result.Data, 2)

	// Verify the result data contains the expected users
	expectedUsers := []string{
		"user1",
		"user2",
	}
	for _, cred := range result.Data {
		require.Equal(t, "svc1", cred.Service)
		require.Contains(t, expectedUsers, cred.Username)
	}

	// Clean up
	testVault.Delete("svc1", "user1")
	testVault.Delete("svc1", "user2")
}

func TestBackupAllCredentials(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/backup",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Add test credentials
	var creds []helpers.FileCred
	for i := range 10 {
		x := rand.Intn(i+1) + 1
		creds = append(creds, helpers.FileCred{
			Service:  fmt.Sprintf("allservice%d", x),
			Username: fmt.Sprintf("alluser%d", i),
			Password: fmt.Sprintf("allpass%d", i),
		})
	}

	t.Cleanup(
		func() {
			for _, cred := range creds {
				err := testVault.Delete(cred.Service, cred.Username)
				if err != nil {
					t.Errorf("cleanup failed for %s/%s: %v", cred.Service, cred.Username, err)
				}
			}
		})

	results := testVault.AddMulti(creds)
	for _, r := range results {
		require.NoError(t, r.Err)
	}

	// Test backing up all credentials
	time.Sleep(2 * time.Second) // wait for eventual consistency
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	cmd := cli.CLICommand{
		File: tmpFile.Name(),
	}
	result := testExecutor.backup(cmd)
	require.True(t, result.Success)
	require.Equal(t, "backup", result.Action)
	require.Len(t, result.Data, len(creds))

	// Verify the result data contains the expected credentials
	expectedCreds := make(map[string][]string)
	for _, cred := range creds {
		expectedCreds[cred.Service] = append(expectedCreds[cred.Service], cred.Username)
	}
	for _, cred := range result.Data {
		expectedUsers, ok := expectedCreds[cred.Service]
		require.True(t, ok, "unexpected credential in result: %s", cred.Service)
		require.Contains(t, expectedUsers, cred.Username)
	}
}
