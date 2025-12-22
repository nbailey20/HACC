package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var ssmPath = "/hackyclient/data/"

func main() {
	fs := flag.NewFlagSet("main", flag.ExitOnError)
	input := parseInput(fs, os.Args[1:]) // Parse the command-line arguments

	// Create backend Vault
	vault, err := NewVault(nil, ssmPath)
	if err != nil {
		fmt.Printf("Error creating vault: %v\n", err)
		return
	}

	// For testing, add some services
	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("service%d", i)
		vault.Add(name, fmt.Sprintf("value%d", i))
	}

	// Pass Vault to Bubbletea model for MVC-like displays
	model := NewModel(input, *vault)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}

	// Clean up test services
	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("service%d", i)
		vault.Delete(name)
	}
}
