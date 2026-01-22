package cli

import (
	"testing"
	"time"

	"github.com/nbailey20/hacc/hacc/config"
	vault "github.com/nbailey20/hacc/hacc/vault"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// Test case 1: hacc
	cmd, err := Parse([]string{}, true)
	require.NoError(t, err)

	_, ok := cmd.Action.(SearchAction)
	require.True(t, ok)
	require.Equal(t, "", cmd.Service)

	// Test case 2: hacc add gmail -u testUser -p secret
	cmd, err = Parse([]string{
		"add",
		"gmail",
		"-u", "testUser",
		"-p", "secret",
	}, true)
	require.NoError(t, err)

	_, ok = cmd.Action.(AddAction)
	require.True(t, ok)

	require.Equal(t, "gmail", cmd.Service)
	require.Equal(t, "testUser", cmd.Username)
	require.Equal(t, "secret", cmd.Password)
	require.False(t, cmd.Generate)

	// Test case 3: hacc a gmail -g
	cmd, err = Parse([]string{
		"a",
		"gmail",
		"-g",
	}, true)
	require.NoError(t, err)

	_, ok = cmd.Action.(AddAction)
	require.True(t, ok)
	require.True(t, cmd.Generate)

	// Test case 4: hacc service1 service2
	_, err = Parse([]string{"service1", "service2"}, true)
	require.Error(t, err)

	// Test case 5: hacc a -g gmail
	cmd, err = Parse([]string{
		"a",
		"-g",
		"gmail",
	}, true)
	require.NoError(t, err)
	require.Equal(t, "gmail", cmd.Service)

	// Test case 6: Positional argument for service at the beginning
	cmd, err = Parse([]string{
		"gmail",
		"a",
		"-g",
	}, true)
	require.NoError(t, err)
	require.Equal(t, "gmail", cmd.Service)

	// Test case 7: Invalid flag (should not set any fields)
	cmd, err = Parse([]string{
		"a",
		"--badflag",
		"gmail",
	}, true)
	require.Error(t, err)

	// Test case 8: hacc add delete -u testUser -g
	// technically works for a service called "delete"
	cmd, err = Parse([]string{
		"add",
		"delete",
		"-u", "testUser",
		"-g",
	}, true)
	require.NoError(t, err)

	// Test case 9: hacc d delete -u testUser
	// technically works for a service called "delete"
	cmd, err = Parse([]string{
		"d",
		"delete",
		"-u", "testUser",
	}, true)
	require.NoError(t, err)

	// Test case 10: hacc search delete
	// technically works for a service called "delete"
	// can't call hacc delete, have to use the implied search keyword
	cmd, err = Parse([]string{
		"search",
		"delete",
	}, true)
	require.NoError(t, err)

	// Test case 11: hacc rotate gmail -u test123 -p 567
	cmd, err = Parse([]string{
		"rotate",
		"gmail",
		"-u", "test123",
		"-p", "567",
	}, true)
	require.NoError(t, err)
}

func TestEqualStringSets(t *testing.T) {
	set1 := []string{"a", "b", "c"}
	set2 := []string{"c", "a", "b"}
	set3 := []string{"a", "b", "d"}
	set4 := []string{"c", "a", "b", "d"}
	if !equalStringSets(set1, set2) {
		t.Errorf("Expected set1 and set2 to be equal")
	}
	if equalStringSets(set1, set3) {
		t.Errorf("Expected set1 and set3 to be not equal")
	}
	if equalStringSets(set1, set4) {
		t.Errorf("Expected set1 and set4 to be not equal")
	}
}

func TestParseAddWithFile(t *testing.T) {
	// Test case 1: hacc add -f backup.json (no service arg)
	cmd, err := Parse([]string{
		"add",
		"-f", "backup.json",
	}, true)
	require.NoError(t, err)

	_, ok := cmd.Action.(AddAction)
	require.True(t, ok)
	require.Equal(t, "backup.json", cmd.File)
	require.Equal(t, "", cmd.Service)

	// Test case 2: hacc a -f backup.json (alias)
	cmd, err = Parse([]string{
		"a",
		"-f", "backup.json",
	}, true)
	require.NoError(t, err)

	_, ok = cmd.Action.(AddAction)
	require.True(t, ok)
	require.Equal(t, "backup.json", cmd.File)
}

func TestValidateCommand(t *testing.T) {
	// Test case 1: Valid add input
	command := CLICommand{
		Action:   AddAction{},
		Service:  "myservice",
		Username: "user1",
		Password: "pass1",
	}
	if err := ValidateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// Test case 2: Valid delete input
	command = CLICommand{Action: DeleteAction{}, Service: "myservice", Username: "user1"}
	if err := ValidateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// Test case 3: Valid search input no username
	command = CLICommand{Action: SearchAction{}, Service: "myservice"}
	if err := ValidateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// Test case 4: Valid search input with username
	command = CLICommand{Action: SearchAction{}, Service: "myservice", Username: "user1"}
	if err := ValidateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// Test case 5: no search input
	command = CLICommand{Action: SearchAction{}}
	if err := ValidateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// Test case 6: Invalid input combination with too many options
	command = CLICommand{Action: AddAction{}, Service: "myservice", Username: "user1", Password: "pass1", Generate: true}
	if err := ValidateCommand(command); err == nil {
		t.Errorf("Expected error for invalid combination, got nil")
	}

	// Test case 7: Invalid input, missing required option
	command = CLICommand{Action: AddAction{}, Username: "user1"}
	if err := ValidateCommand(command); err == nil {
		t.Errorf("Expected error for missing required option, got nil")
	}
}

func TestAutoCompleteCommand(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath: "/hackyclient/test/",
	}
	vault, err := vault.NewVault(nil, cfg)
	if err != nil {
		t.Fatalf("Error creating vault: %v", err)
	}
	// Prepopulate vault with services
	err = vault.Add("serv1Name", "user1", "value1")
	if err != nil {
		t.Fatalf("Error adding service1: %v", err)
	}
	err = vault.Add("serv1Name", "user1a", "value1a")
	if err != nil {
		t.Fatalf("Error adding 2nd user to service1: %v", err)
	}
	err = vault.Add("serv2Name", "user2", "value2")
	if err != nil {
		t.Fatalf("Error adding service2: %v", err)
	}

	// Test case 1: Exact service match
	command := CLICommand{Service: "serv1Name", Action: SearchAction{}}
	updatedCommand := AutoCompleteCommand(command, vault)
	if updatedCommand.Service != "serv1Name" {
		t.Errorf("Expected service to be 'serv1Name', got %v", updatedCommand.Service)
	}

	// Test case 2: Single partial service match
	command = CLICommand{Service: "serv1", Action: SearchAction{}}
	updatedCommand = AutoCompleteCommand(command, vault)
	if updatedCommand.Service != "serv1Name" {
		t.Errorf("Expected service to be autocompleted to 'serv1Name', got %v", updatedCommand.Service)
	}

	// Test case 3: Multiple partial service matches
	command = CLICommand{Service: "serv", Action: SearchAction{}}
	updatedCommand = AutoCompleteCommand(command, vault)
	if updatedCommand.Service != "serv" {
		t.Errorf("Expected service to remain 'serv' due to multiple matches, got %v", updatedCommand.Service)
	}

	// Test case 4: No service match
	command = CLICommand{Service: "unknown", Action: SearchAction{}}
	updatedCommand = AutoCompleteCommand(command, vault)
	if updatedCommand.Service != "unknown" {
		t.Errorf("Expected service to remain 'unknown' due to no match, got %v", updatedCommand.Service)
	}

	// Test case 5: Single partial username match
	command = CLICommand{Service: "serv2Name", Username: "u", Action: SearchAction{}}
	updatedCommand = AutoCompleteCommand(command, vault)
	if updatedCommand.Username != "user2" {
		t.Errorf("Expected username to be autocompleted to 'user2', got %v", updatedCommand.Username)
	}

	// Test case 6: Multiple partial username matches
	command = CLICommand{Service: "serv1Name", Username: "user1", Action: SearchAction{}}
	updatedCommand = AutoCompleteCommand(command, vault)
	if updatedCommand.Username != "user1" {
		t.Errorf("Expected username to remain 'user1' due to multiple matches, got %v", updatedCommand.Username)
	}

	// Test case 7: No username match
	command = CLICommand{Service: "serv1Name", Username: "unknown", Action: SearchAction{}}
	updatedCommand = AutoCompleteCommand(command, vault)
	if updatedCommand.Username != "unknown" {
		t.Errorf("Expected username to remain 'unknown' due to no match, got %v", updatedCommand.Username)
	}

	// Clean up
	err = vault.Delete("serv1Name", "user1")
	if err != nil {
		t.Fatalf("Error deleting service1: %v", err)
	}
	err = vault.Delete("serv1Name", "user1a")
	if err != nil {
		t.Fatalf("Error deleting 2nd user from service1: %v", err)
	}
	err = vault.Delete("serv2Name", "user2")
	if err != nil {
		t.Fatalf("Error deleting service2: %v", err)
	}
	time.Sleep(2 * time.Second) // wait for eventual consistency
}
