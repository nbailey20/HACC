package engine

import (
	"os"
	"testing"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/config"
	"github.com/nbailey20/hacc/vault"
	"github.com/stretchr/testify/require"
)

func TestDeleteCredential(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/delete",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Add a test credential
	err = testVault.Add("deleteservice", "deleteuser", "deletepass")
	require.NoError(t, err)

	// Test deleting the credential
	cmd := cli.CLICommand{
		Action:   cli.DeleteAction{},
		Service:  "deleteservice",
		Username: "deleteuser",
	}
	result := testExecutor.Execute(cmd)
	require.True(t, result.Success)
	require.Equal(t, "delete", result.Action)
	require.Len(t, result.Data, 1)
	require.Equal(t, "deleteservice", result.Data[0].Service)
	require.Equal(t, "deleteuser", result.Data[0].Username)
}

func TestDeleteMultiCredential(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/delete",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Add multiple test credentials to vault
	err = testVault.Add("deleteservice1", "deleteuser1", "deletepass1")
	require.NoError(t, err)
	err = testVault.Add("deleteservice2", "deleteuser2", "deletepass2")
	require.NoError(t, err)

	// Create a temporary JSON file with test credentials
	tmpFile, err := os.CreateTemp("", "delete*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	testData := `{"creds_list":[{"service":"deleteservice1","username":"deleteuser1"},{"service":"deleteservice2","username":"deleteuser2"}]}`
	_, err = tmpFile.WriteString(testData)
	require.NoError(t, err)
	tmpFile.Close()

	// Test deleting multiple credentials
	cmd := cli.CLICommand{
		Action: cli.DeleteAction{},
		File:   tmpFile.Name(),
	}
	result := testExecutor.Execute(cmd)
	require.True(t, result.Success)
	require.Equal(t, "delete", result.Action)
	require.Len(t, result.Data, 2)
	services := []string{result.Data[0].Service, result.Data[1].Service}
	usernames := []string{result.Data[0].Username, result.Data[1].Username}
	require.Contains(t, services, "deleteservice1")
	require.Contains(t, services, "deleteservice2")
	require.Contains(t, usernames, "deleteuser1")
	require.Contains(t, usernames, "deleteuser2")
}
