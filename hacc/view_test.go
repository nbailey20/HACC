package main

import (
	"fmt"
	"strings"
	"testing"

	vault "github.com/nbailey20/hacc/hacc/vault"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDisplay(t *testing.T) {
	testVault, err := vault.NewVault(nil, "/hackyclient/test/display")
	if err != nil {
		t.Fatalf("Error creating vault: %v", err)
	}

	model := model{
		vault:      *testVault,
		pageSize:   2,
		showSecret: false,
		state:      &ServiceListState{},
	}

	// Add services to the vault
	for i := 1; i <= 5; i++ {
		serviceName := fmt.Sprintf("service%d", i)
		userName := fmt.Sprintf("user%d", i)
		value := fmt.Sprintf("value%d", i)
		err := testVault.Add(serviceName, userName, value)
		if err != nil {
			t.Fatalf("Error adding service: %v", err)
		}
	}

	// Simulate displaying the list view
	view := model.ServiceListView()
	expectedLines := []string{
		fmt.Sprintf("%s %-5d %-25s", ">", 1, "service1"),
		fmt.Sprintf("%s %-5d %-25s", " ", 2, "service2"),
	}
	for _, line := range expectedLines {
		if !strings.Contains(view, line) {
			t.Errorf("Expected line '%s' in view, got:\n%s", line, view)
		}
	}

	// Simulate navigating down
	model.cursor = 1
	view = model.ServiceListView()
	if !strings.Contains(view, fmt.Sprintf("%s %-5d %-25s", ">", 2, "service2")) {
		t.Errorf("Expected cursor on service2, got:\n%s", view)
	}
	// Simulate page change
	model.page = 1
	view = model.ServiceListView()
	if !strings.Contains(view, fmt.Sprintf("%s %-5d %-25s", ">", 2, "service4")) {
		t.Errorf("Expected cursor on service4 on page 2, got:\n%s", view)
	}

	// Delete the services after test
	for i := 1; i <= 5; i++ {
		serviceName := fmt.Sprintf("service%d", i)
		userName := fmt.Sprintf("user%d", i)
		err := testVault.Delete(serviceName, userName)
		if err != nil {
			t.Fatalf("Error deleting service: %v", err)
		}
	}
	// Simulate quitting
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatalf("expected a command, got nil")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", msg)
	}
}
