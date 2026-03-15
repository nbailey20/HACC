package engine

import (
	"testing"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/config"
	"github.com/nbailey20/hacc/vault"
	"github.com/stretchr/testify/require"
)

func TestRotateCredential(t *testing.T) {
	cfg := config.Config{
		ParamPath:      "/hackyclient/test/rotate",
		ObfuscationKey: "secretkey",
	}
	// Create a temporary vault / executor
	testVault, err := vault.NewVault(nil, &cfg)
	require.NoError(t, err)
	testExecutor := NewExecutor(testVault)

	// Add a test credential
	err = testVault.Add("rotateservice", "rotateuser", "rotatepass")
	require.NoError(t, err)

	// Test rotating the credential
	cmd := cli.CLICommand{
		Action:   cli.RotateAction{},
		Service:  "rotateservice",
		Username: "rotateuser",
		Password: "newrotatepass",
	}
	result := testExecutor.Execute(cmd)
	require.True(t, result.Success)
	require.Equal(t, "rotate", result.Action)
	require.Len(t, result.Data, 1)
	require.Equal(t, "rotateservice", result.Data[0].Service)
	require.Equal(t, "rotateuser", result.Data[0].Username)

	// Clean up
	testVault.Delete("rotateservice", "rotateuser")
}
