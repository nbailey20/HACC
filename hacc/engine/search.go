package engine

import (
	"fmt"

	"github.com/nbailey20/hacc/cli"
)

func (e *Executor) search(cmd cli.CLICommand) Result {
	services := e.vault.ListServices(cmd.Service)
	if len(services) == 0 {
		return Result{
			Action:  "search",
			Success: false,
			Error:   "no matching services found",
		}
	}
	if len(services) > 1 && cmd.Username != "" {
		return Result{
			Action:  "search",
			Success: false,
			Error:   fmt.Sprintf("ambiguous search: multiple services match '%s': %v", cmd.Service, services),
		}
	}
	// if multiple services match and no username provided, return list of matching services
	if len(services) > 1 {
		var creds []Credential
		for _, svc := range services {
			creds = append(creds, Credential{Service: svc})
		}
		return Result{
			Action:  "search",
			Success: true,
			Data:    creds,
		}
	}
	// if only one service matches but no username provided, return list of users for that service
	if len(services) == 1 && cmd.Username == "" {
		users, err := e.vault.GetUsersForService(services[0])
		if err != nil {
			return Result{
				Action:  "search",
				Success: false,
				Error:   err.Error(),
			}
		}
		var creds []Credential
		for _, user := range users {
			creds = append(creds, Credential{Service: services[0], Username: user})
		}
		return Result{
			Action:  "search",
			Success: true,
			Data:    creds,
		}
	}
	// if only one service matches and username provided, return the credential
	password, err := e.vault.Get(services[0], cmd.Username)
	if err != nil {
		return Result{
			Action:  "search",
			Success: false,
			Error:   err.Error(),
		}
	}
	return Result{
		Action:  "search",
		Success: true,
		Data:    []Credential{{Service: services[0], Username: cmd.Username, Password: password}},
	}
}
