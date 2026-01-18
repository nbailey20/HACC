package display

import (
	"fmt"
	"os"
	"strings"
	"testing"

	vault "github.com/nbailey20/hacc/hacc/vault"
	"github.com/stretchr/testify/require"

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

func TestAddMultiCredential(t *testing.T) {
	// Create a temporary vault
	testVault, err := vault.NewVault(nil, "/hackyclient/test/multi/")
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
	addMultiCmd := addMultiCredentialCmd(*testVault, tmpFile.Name())
	msg := addMultiCmd()

	// Should return AddFailedMsg with no error if all succeed
	failMsg, ok := msg.(AddFailedMsg)
	require.True(t, ok)
	require.NoError(t, failMsg.Error)

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
	testVault, err := vault.NewVault(nil, "/hackyclient/test/multi/")
	require.NoError(t, err)

	// Test case: File doesn't exist
	addMultiCmd := addMultiCredentialCmd(*testVault, "nonexistent.json")
	msg := addMultiCmd()

	failMsg, ok := msg.(AddFailedMsg)
	require.True(t, ok)
	require.Error(t, failMsg.Error)
}

func TestAddMultiCredentialInvalidJson(t *testing.T) {
	testVault, err := vault.NewVault(nil, "/hackyclient/test/multi/")
	require.NoError(t, err)

	// Create a temporary file with invalid JSON
	tmpFile, err := os.CreateTemp("", "backup*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("invalid json {")
	tmpFile.Close()

	// Test case: Invalid JSON
	addMultiCmd := addMultiCredentialCmd(*testVault, tmpFile.Name())
	msg := addMultiCmd()

	failMsg, ok := msg.(AddFailedMsg)
	require.True(t, ok)
	require.Error(t, failMsg.Error)
}
