package main

import (
	"fmt"

	cli "github.com/nbailey20/hacc/hacc/cli"
	vault "github.com/nbailey20/hacc/hacc/vault"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	vault       vault.Vault
	serviceName string
	username    string
	page        int
	pageSize    int
	cursor      int
	message     string
	showSecret  bool
	state       State
}

func (m model) Init() tea.Cmd {
	return nil
}

func NewModel(cli cli.CLICommand, vault vault.Vault) *model {
	switch cli.Action {
	case "search":
		if cli.Service != "" {
			return &model{
				vault:       vault,
				serviceName: cli.Service,
				username:    cli.Username,
				page:        0,
				pageSize:    pageSize,
				cursor:      0,
				message:     "",
				showSecret:  true,
				state:       &CredentialState{},
			}
		}
		return &model{
			vault:       vault,
			serviceName: "",
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  true,
			state:       &WelcomeState{},
		}
	case "add":
		vault.Add(cli.Service, cli.Username, cli.Password)
		return &model{
			vault:       vault,
			serviceName: cli.Service,
			username:    cli.Username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     fmt.Sprintf("Successfully saved credential for %s", cli.Service),
			showSecret:  false,
			state:       &EndState{},
		}
	case "delete":
		err := vault.Delete(cli.Service, cli.Username)
		msg := fmt.Sprintf("Service %s deleted successfully.", cli.Service)
		if err != nil {
			msg = fmt.Sprintf("Error deleting username %s in service %s: %v\n", cli.Username, cli.Service, err)
		}
		return &model{
			vault:       vault,
			serviceName: "",
			username:    cli.Username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     msg,
			showSecret:  false,
			state:       &EndState{},
		}
	default:
		return &model{
			vault:       vault,
			serviceName: "",
			username:    cli.Username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  false,
			state:       &WelcomeState{},
		}
	}
}
