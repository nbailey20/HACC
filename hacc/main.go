package main

import (
	"fmt"
	"os"

	cli "github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/config"
	display "github.com/nbailey20/hacc/display"
	"github.com/nbailey20/hacc/engine"
	vault "github.com/nbailey20/hacc/vault"
)

func main() {
	// Determine config.yaml location
	path, err := config.GetPath()
	if err != nil {
		fmt.Printf("Error determining path of config.yaml, exiting: %v\n", err)
		return
	}

	// Load config data
	cfg, err := config.Load(path)
	if err != nil {
		fmt.Printf("Error loading contents of %s, exiting: %v\n", path, err)
		return
	}
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Error loading config file, exiting: %v\n", err)
		return
	}

	// Parse the command-line arguments
	command, err := cli.Parse(os.Args[1:], false)
	if err != nil {
		fmt.Printf("Error parsing input, exiting: %v\n", err)
		return
	}

	// Create backend Vault, don't pass any services
	// to auto-load from backend
	vault, err := vault.NewVault(nil, cfg)
	if err != nil {
		fmt.Printf("Error creating vault, exiting: %v\n", err)
		return
	}

	// Auto-complete command where possible given vault data
	// purely for enhancing UX and less typing
	command = cli.AutoCompleteCommand(command, vault)

	// Pass Vault executor to display component and start program
	executor := engine.NewExecutor(vault)
	if err := display.Start(command, executor); err != nil {
		fmt.Println("Error starting program:", err)
	}
}
