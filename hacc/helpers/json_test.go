package helpers

import (
	"os"
	"testing"
)

func TestReadJsonFile(t *testing.T) {
	// Create a temporary JSON file
	tmpFile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test JSON data
	testData := `{"creds_list":[{"service":"github","username":"user1","password":"pass1"}]}`
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Test successful read
	result, err := ReadJsonFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	creds := result["creds_list"]
	if len(creds) == 0 {
		t.Error("Expected creds_list to have at least one entry")
	}

	if creds[0]["service"] != "github" {
		t.Errorf("Expected service 'github', got: %s", creds[0]["service"])
	}
}

func TestReadJsonFileNotFound(t *testing.T) {
	_, err := ReadJsonFile("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestReadJsonFileInvalid(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("invalid json {")
	tmpFile.Close()

	_, err = ReadJsonFile(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
