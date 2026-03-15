package engine

import (
	"testing"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/config"
	"github.com/nbailey20/hacc/vault"
	"github.com/stretchr/testify/require"
)

func TestSearchCredential(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/search",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Add a test credential
	err = testVault.Add("searchservice", "searchuser", "searchpass")
	require.NoError(t, err)

	// Test searching for the credential
	cmd := cli.CLICommand{
		Action:   cli.SearchAction{},
		Service:  "searchservice",
		Username: "searchuser",
	}
	result := testExecutor.Execute(cmd)
	require.True(t, result.Success)
	require.Equal(t, "search", result.Action)
	require.Len(t, result.Data, 1)
	require.Equal(t, "searchservice", result.Data[0].Service)
	require.Equal(t, "searchuser", result.Data[0].Username)
	require.Equal(t, "searchpass", result.Data[0].Password)

	// Clean up
	testVault.Delete("searchservice", "searchuser")
}
