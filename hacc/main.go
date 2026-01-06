package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var ssmPath = "/hackyclient/data/"

func main() {
	// Parse the command-line arguments
	fs := flag.NewFlagSet("main", flag.ExitOnError)
	input, err := parseInput(fs, os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}

	// Normalize input into raw CLI command
	rawCli := normalizeInput(input)

	// Validate raw command into final CLI command
	cli, err := validateCommand(rawCli)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	// Create backend Vault
	vault, err := NewVault(nil, ssmPath)
	if err != nil {
		fmt.Printf("Error creating vault: %v\n", err)
		return
	}

	// Auto-complete command where possible given vault data
	// purely for enhancing UX and less typing
	cli = autoCompleteCommand(cli, vault)

	// Pass Vault to Bubbletea model for MVC-like displays (interactive mode)
	model := NewModel(cli, *vault)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}
}
