package cli

import (
	"github.com/nbailey20/hacc/vault"
)

func AutoCompleteCommand(command CLICommand, vault *vault.Vault) CLICommand {
	// Auto-complete missing fields in commands where possible using vault data
	command = autoCompleteService(command, vault)
	command = autoCompleteUsername(command, vault)
	return command
}
func autoCompleteService(command CLICommand, vault *vault.Vault) CLICommand {
	// If only one service exists in vault for a given prefix, auto-complete it
	// Do not autocomplete any Add actions
	if command.Action.Kind() != ActionAdd && len(vault.ListServices(command.Service)) == 1 {
		command.Service = vault.ListServices(command.Service)[0]
	}
	return command
}

func autoCompleteUsername(command CLICommand, vault *vault.Vault) CLICommand {
	// If only one user prefix exists for a valid service, auto-complete it
	if command.Action.Kind() != ActionAdd && vault.HasService(command.Service) {
		if len(vault.Services[command.Service].GetUsers(command.Username)) == 1 {
			command.Username = vault.Services[command.Service].GetUsers(command.Username)[0]
		}
	}
	return command
}
