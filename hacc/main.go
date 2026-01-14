package main

import (
	"fmt"
	"os"

	cli "github.com/nbailey20/hacc/hacc/cli"
	display "github.com/nbailey20/hacc/hacc/display"
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
		// If we received a help command for usage, we're done
		if err.Error() == "No action" {
			return
		}
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

	// Pass Vault to display component and start program
	if err := display.Start(command, vault); err != nil {
		fmt.Println("Error starting program:", err)
	}
}
