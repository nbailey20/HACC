package main

import (
	"flag"
	"testing"
)

func TestParseInput(t *testing.T) {
	// Test case 1: Default values
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	input := parseInput(fs, []string{})
	if input.Action != "get" {
		t.Errorf("Expected default action 'get', got '%s'", input.Action)
	}
	if input.Name != "" {
		t.Errorf("Expected default name '', got '%s'", input.Name)
	}
	if input.Value != "" {
		t.Errorf("Expected default value '', got '%s'", input.Value)
	}

	// Test case 2: Custom values
	fs = flag.NewFlagSet("test2", flag.ExitOnError)
	args := []string{"-action", "add", "-name", "testService", "-value", "testValue"}
	input = parseInput(fs, args)
	if input.Action != "add" {
		t.Errorf("Expected action 'add', got '%s'", input.Action)
	}
	if input.Name != "testService" {
		t.Errorf("Expected name 'testService', got '%s'", input.Name)
	}
	if input.Value != "testValue" {
		t.Errorf("Expected value 'testValue', got '%s'", input.Value)
	}
}
