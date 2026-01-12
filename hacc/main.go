package main

import (
	"fmt"
	"os"

	cli "github.com/nbailey20/hacc/hacc/cli"
	vault "github.com/nbailey20/hacc/hacc/vault"
)

var ssmPath = "/hackyclient/data/"

func main() {
	// Parse the command-line arguments
	command, err := cli.Parse(os.Args[1:], false)
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}

	// Validate raw command into final CLI command
	err = cli.ValidateCommand(command)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	// Create backend Vault
	vault, err := vault.NewVault(nil, ssmPath)
	if err != nil {
		fmt.Printf("Error creating vault: %v\n", err)
		return
	}

	// Auto-complete command where possible given vault data
	// purely for enhancing UX and less typing
	command = cli.AutoCompleteCommand(command, vault)
	fmt.Println("Final CLI Command:", command)

	// // Pass Vault to Bubbletea model for MVC-like displays (interactive mode)
	// model := NewModel(command, *vault)
	// p := tea.NewProgram(model)
	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("Error starting program:", err)
	// }
}
