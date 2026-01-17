package display

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

	// Check for expected footer
	view := model.ServiceListView()
	expectedLines := []string{
		"│ # │ Service  │",
		"│ 1 │service1  │",
		"│ 2 │ service2 │",
		"1 / 3",
		"←→↑↓ / 1-9: navigate   Enter: select",
		"Backspace: back        Esc / ^C: quit",
	}
	for _, line := range expectedLines {
		if !strings.Contains(view, line) {
			t.Errorf("Expected line '%s' in view, got:\n%s", line, view)
		}
	}

	// Simulate navigating down
	model.cursor = 1
	view = model.ServiceListView()
	expectedLines = []string{
		"│ 1 │ service1 │",
		"│ 2 │service2  │",
		"1 / 3",
	}
	for _, line := range expectedLines {
		if !strings.Contains(view, line) {
			t.Errorf("Expected line '%s' in view, got:\n%s", line, view)
		}
	}
	// Simulate page change
	model.page = 1
	view = model.ServiceListView()
	expectedLines = []string{
		"│ 1 │ service3 │",
		"│ 2 │service4  │",
		"2 / 3",
	}
	for _, line := range expectedLines {
		if !strings.Contains(view, line) {
			t.Errorf("Expected line '%s' in view, got:\n%s", line, view)
		}
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
