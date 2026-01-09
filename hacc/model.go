package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	vault       Vault
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

func NewModel(cli CLICommand, vault Vault) *model {
	switch cli.action {
	case "search":
		if cli.service != "" {
			return &model{
				vault:       vault,
				serviceName: cli.service,
				username:    cli.username,
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
		vault.Add(cli.service, cli.username, cli.password)
		return &model{
			vault:       vault,
			serviceName: cli.service,
			username:    cli.username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     fmt.Sprintf("Successfully saved credential for %s", cli.service),
			showSecret:  false,
			state:       &EndState{},
		}
	case "delete":
		err := vault.Delete(cli.service, cli.username)
		msg := fmt.Sprintf("Service %s deleted successfully.", cli.service)
		if err != nil {
			msg = fmt.Sprintf("Error deleting username %s in service %s: %v\n", cli.username, cli.service, err)
		}
		return &model{
			vault:       vault,
			serviceName: "",
			username:    cli.username,
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
			username:    cli.username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  false,
			state:       &WelcomeState{},
		}
	}
}
