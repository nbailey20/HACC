package display

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nbailey20/hacc/hacc/cli"
	"github.com/nbailey20/hacc/hacc/config"
	"github.com/nbailey20/hacc/hacc/helpers"
	vault "github.com/nbailey20/hacc/hacc/vault"
)

func TestAddMultiCredential(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath:      "/hackyclient/test/multi",
		ObfuscationKey: "secretkey",
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
	addMultiCmd := addMultiCredentialCmd(testVault, tmpFile.Name())
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
		ParamPath:      "/hackyclient/test/multi",
		ObfuscationKey: "secretkey",
	}
	testVault, err := vault.NewVault(nil, cfg)
	require.NoError(t, err)

	// Test case: File doesn't exist
	addMultiCmd := addMultiCredentialCmd(testVault, "nonexistent.json")
	msg := addMultiCmd()

	failMsg, ok := msg.(AddFailedMsg)
	require.True(t, ok)
	require.Error(t, failMsg.Error)
}

func TestAddMultiCredentialInvalidJson(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath:      "/hackyclient/test/multi",
		ObfuscationKey: "secretkey",
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
	addMultiCmd := addMultiCredentialCmd(testVault, tmpFile.Name())
	msg := addMultiCmd()

	failMsg, ok := msg.(AddFailedMsg)
	require.True(t, ok)
	require.Error(t, failMsg.Error)
}

func TestBackupCredential(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath:      "/hackyclient/test/backup",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault
	testVault, err := vault.NewVault(nil, cfg)
	require.NoError(t, err)
	// Add a test credential
	err = testVault.Add("backupservice", "backupuser", "backuppass")
	require.NoError(t, err)

	// Test backing up the credential to a file
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	cmd := cli.CLICommand{
		File:     tmpFile.Name(),
		Service:  "backupservice",
		Username: "backupuser",
	}
	msg := backupCmd(testVault, cmd)()

	backupMsg, ok := msg.(BackupSuccessMsg)
	require.True(t, ok)
	require.Contains(t, backupMsg.Display, "backed up credential backupservice/backupuser")
	// Clean up
	testVault.Delete("backupservice", "backupuser")
}

func TestBackupServiceCredentials(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath:      "/hackyclient/test/backup",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault
	testVault, err := vault.NewVault(nil, cfg)
	require.NoError(t, err)
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
	msg := backupCmd(testVault, cmd)()
	backupMsg, ok := msg.(BackupSuccessMsg)
	require.True(t, ok)
	require.Contains(t, backupMsg.Display, "backed up credential svc1/user1")
	require.Contains(t, backupMsg.Display, "backed up credential svc1/user2")
	// Clean up
	testVault.Delete("svc1", "user1")
	testVault.Delete("svc1", "user2")
}

func TestBackupAllCredentials(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath:      "/hackyclient/test/backup",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault
	testVault, err := vault.NewVault(nil, cfg)
	require.NoError(t, err)
	// Add test credentials
	var creds []helpers.FileCred
	for i := range 20 {
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
	msg := backupCmd(testVault, cmd)()
	backupMsg, ok := msg.(BackupSuccessMsg)
	require.True(t, ok)
	for _, cred := range creds {
		require.Contains(t, backupMsg.Display, "backed up credential "+cred.Service+"/"+cred.Username)
	}
}
