package main

import (
	"flag"
	"io"
	"testing"
	"time"
)

func TestParseInput(t *testing.T) {
	// Test case 1: Default values
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	input, err := parseInput(fs, []string{})
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}
	if input.username != "" {
		t.Errorf("Expected username to be empty, got %v", input.username)
	}

	// Test case 2: Set install flag
	fs = flag.NewFlagSet("testi", flag.ExitOnError)
	input, err = parseInput(fs, []string{"--install"})
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}
	// Test case 3: Set username and password
	fs = flag.NewFlagSet("testu", flag.ExitOnError)
	input, err = parseInput(fs, []string{"--username", "testuser", "--password", "testpass"})
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}
	if input.username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %v", input.username)
	}
	if input.password != "testpass" {
		t.Errorf("Expected password to be 'testpass', got %v", input.password)
	}

	// Test case 4: Set multiple flags
	fs = flag.NewFlagSet("testm", flag.ExitOnError)
	input, err = parseInput(fs, []string{"--add", "--rotate", "--file", "config.yaml"})
	if input.add != true {
		t.Errorf("Expected add to be true, got %v", input.add)
	}
	if input.rotate != true {
		t.Errorf("Expected rotate to be true, got %v", input.rotate)
	}
	if input.file != "config.yaml" {
		t.Errorf("Expected file to be 'config.yaml', got %v", input.file)
	}

	// Test case 5: Positional argument for service at the end
	fs = flag.NewFlagSet("tests", flag.ExitOnError)
	input, err = parseInput(fs, []string{"--add", "--username", "user1", "--password", "pass1", "myservice"})
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}
	if input.service != "myservice" {
		t.Errorf("Expected service to be 'myservice', got %v", input.service)
	}
	if input.username != "user1" {
		t.Errorf("Expected username to be 'user1', got %v", input.username)
	}
	if input.password != "pass1" {
		t.Errorf("Expected password to be 'pass1', got %v", input.password)
	}

	// // Test case 6: Positional argument for service at the beginning
	// flags library doesn't support positional args at the beginning or middle
	// fs = flag.NewFlagSet("tests", flag.ExitOnError)
	// input, err = parseInput(fs, []string{"myservice", "--add", "--username", "user1", "--password", "pass1"})
	// if err != nil {
	// 	t.Errorf("Error parsing input: %v", err)
	// }
	// if input.service != "myservice" {
	// 	t.Errorf("Expected service to be 'myservice', got %v", input.service)
	// }
	// if input.username != "user1" {
	// 	t.Errorf("Expected username to be 'user1', got %v", input.username)
	// }
	// if input.password != "pass1" {
	// 	t.Errorf("Expected password to be 'pass1', got %v", input.password)
	// }

	// Test case 6: Invalid flag (should not set any fields)
	fs = flag.NewFlagSet("testin", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input, err = parseInput(fs, []string{"--invalidflag"})
	if err == nil {
		t.Errorf("Expected error for invalid flag, got nil")
	}
	if input.username != "" {
		t.Errorf("Expected username to be empty for invalid flag, got %v", input.username)
	}
}

func TestNormalizeInput(t *testing.T) {
	// Test case 1: Valid Add action
	input := CLIInput{add: true, service: "myservice", username: "user1", password: "pass1"}
	normalized := normalizeInput(input)
	if normalized.action[0] != "add" {
		t.Errorf("Expected action[0] to be 'add', got %v", normalized.action[0])
	}
	if normalized.service != "myservice" {
		t.Errorf("Expected service to be 'myservice', got %v", normalized.service)
	}
	if normalized.username != "user1" {
		t.Errorf("Expected username to be 'user1', got %v", normalized.username)
	}
	if normalized.password != "pass1" {
		t.Errorf("Expected password to be 'pass1', got %v", normalized.password)
	}

	// Test case 2: Valid Delete action
	input = CLIInput{delete: true, service: "myservice", username: "user1"}
	normalized = normalizeInput(input)
	if normalized.action[0] != "delete" {
		t.Errorf("Expected action to be 'delete', got %v", normalized.action[0])
	}
	if normalized.service != "myservice" {
		t.Errorf("Expected service to be 'myservice', got %v", normalized.service)
	}
	if normalized.username != "user1" {
		t.Errorf("Expected username to be 'user1', got %v", normalized.username)
	}

	// Test case 3: Default Search action
	input = CLIInput{service: "myservice"}
	normalized = normalizeInput(input)
	if normalized.action[0] != "search" {
		t.Errorf("Expected action to be 'search', got %v", normalized.action[0])
	}
	if normalized.service != "myservice" {
		t.Errorf("Expected service to be 'myservice', got %v", normalized.service)
	}

	// Test case 4: Multiple action flags which will later fail validation (lossless normalization)
	input = CLIInput{add: true, delete: true, service: "myservice", username: "user1"}
	normalized = normalizeInput(input)
	if len(normalized.action) != 2 {
		t.Errorf("Expected action to be [add, delete], got %v", normalized.action)
	}
	if normalized.service != "myservice" {
		t.Errorf("Expected service to be 'myservice', got %v", normalized.service)
	}
	if normalized.username != "user1" {
		t.Errorf("Expected username to be 'user1', got %v", normalized.username)
	}
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

func TestValidateCommand(t *testing.T) {
	// Test case 1: Valid add input
	command := RawCLICommand{
		action:   []string{"add"},
		service:  "myservice",
		username: "user1",
		password: "pass1",
	}
	if _, err := validateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// Test case 2: Valid delete input
	command = RawCLICommand{action: []string{"delete"}, service: "myservice", username: "user1"}
	if _, err := validateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// // Test case 3: Valid search input no username
	// input = CLIInput{service: "myservice"}
	// if err := validateCommand(input); err != nil {
	// 	t.Errorf("Expected no error for valid input, got %v", err)
	// }

	// Test case 4: Valid search input with username
	command = RawCLICommand{action: []string{"search"}, service: "myservice", username: "user1"}
	if _, err := validateCommand(command); err != nil {
		t.Errorf("Expected no error for valid input, got %v", err)
	}

	// Test case 5: Invalid input, multiple actions not allowed
	command = RawCLICommand{action: []string{"add", "delete"}, service: "myservice", username: "user1"}
	if _, err := validateCommand(command); err == nil {
		t.Errorf("Expected error for multiple actions, got nil")
	}

	// Test case 6: Invalid input combination with too many options
	command = RawCLICommand{action: []string{"add"}, service: "myservice", username: "user1", password: "pass1", generate: true}
	if _, err := validateCommand(command); err == nil {
		t.Errorf("Expected error for invalid combination, got nil")
	}

	// Test case 7: Invalid input, missing required option
	command = RawCLICommand{action: []string{"add"}, username: "user1"}
	if _, err := validateCommand(command); err == nil {
		t.Errorf("Expected error for missing required option, got nil")
	}
}

func TestAutoCompleteCommand(t *testing.T) {
	vault, err := NewVault(nil, "/hackyclient/test/")
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
	command := CLICommand{service: "serv1Name"}
	updatedCommand := autoCompleteCommand(command, vault)
	if updatedCommand.service != "serv1Name" {
		t.Errorf("Expected service to be 'serv1Name', got %v", updatedCommand.service)
	}

	// Test case 2: Single partial service match
	command = CLICommand{service: "serv1"}
	updatedCommand = autoCompleteCommand(command, vault)
	if updatedCommand.service != "serv1Name" {
		t.Errorf("Expected service to be autocompleted to 'serv1Name', got %v", updatedCommand.service)
	}

	// Test case 3: Multiple partial service matches
	command = CLICommand{service: "serv"}
	updatedCommand = autoCompleteCommand(command, vault)
	if updatedCommand.service != "serv" {
		t.Errorf("Expected service to remain 'serv' due to multiple matches, got %v", updatedCommand.service)
	}

	// Test case 4: No service match
	command = CLICommand{service: "unknown"}
	updatedCommand = autoCompleteCommand(command, vault)
	if updatedCommand.service != "unknown" {
		t.Errorf("Expected service to remain 'unknown' due to no match, got %v", updatedCommand.service)
	}

	// Test case 5: Single partial username match
	command = CLICommand{service: "serv2Name", username: "u"}
	updatedCommand = autoCompleteCommand(command, vault)
	if updatedCommand.username != "user2" {
		t.Errorf("Expected username to be autocompleted to 'user2', got %v", updatedCommand.username)
	}

	// Test case 6: Multiple partial username matches
	command = CLICommand{service: "serv1Name", username: "user1"}
	updatedCommand = autoCompleteCommand(command, vault)
	if updatedCommand.username != "user1" {
		t.Errorf("Expected username to remain 'user1' due to multiple matches, got %v", updatedCommand.username)
	}

	// Test case 7: No username match
	command = CLICommand{service: "serv1Name", username: "unknown"}
	updatedCommand = autoCompleteCommand(command, vault)
	if updatedCommand.username != "unknown" {
		t.Errorf("Expected username to remain 'unknown' due to no match, got %v", updatedCommand.username)
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
