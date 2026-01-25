package helpers

import (
	"os"
	"testing"
)

func TestReadWriteJsonFile(t *testing.T) {
	// Create a temporary JSON file
	tmpFile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test JSON data
	testData := map[string][]map[string]string{
		"creds_list": {{"service": "github", "username": "user1", "password": "pass1"}},
	}
	if err := writeJsonFile(tmpFile.Name(), testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Test successful read
	result, err := readJsonFile(tmpFile.Name())
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
	_, err := readJsonFile("nonexistent.json")
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

	_, err = readJsonFile(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestReadWriteCredsFile(t *testing.T) {
	// Create a temporary JSON file
	tmpFile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test JSON data
	testData := []FileCred{{Service: "github", Username: "user1", Password: "pass1"}}
	if err := WriteCredsFile(tmpFile.Name(), testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Test successful read
	result, err := ReadCredsFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 1 {
		t.Error("Expected one FileCred result")
	}

	if result[0].Service != "github" {
		t.Errorf("Expected service 'github', got: %s", result[0].Service)
	}
	if result[0].Username != "user1" {
		t.Errorf("Expected username 'user1', got: %s", result[0].Username)
	}
	if result[0].Password != "pass1" {
		t.Errorf("Expected password 'pass1', got: %s", result[0].Password)
	}
}
