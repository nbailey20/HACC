package display

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

func searchModelState(service string, user string, vault vault.Vault) State {
	// determines starting state for search commands
	// if we already have enough info from CLI/autocompletion,
	// jump straight to credential view
	if service == "" && user == "" {
		return &WelcomeState{}
	}
	if vault.HasService(service) && vault.Services[service].HasUser(user) {
		return &CredentialState{}
	}
	if vault.HasService(service) {
		return &UsernameListState{}
	}
	return &ServiceListState{}
}

func NewModel(cmd cli.CLICommand, vault vault.Vault) *model {
	switch cmd.Action.(type) {
	case cli.SearchAction:
		return &model{
			vault:       vault,
			serviceName: cmd.Service,
			username:    cmd.Username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  true,
			state:       searchModelState(cmd.Service, cmd.Username, vault),
		}
	case cli.AddAction:
		err := vault.Add(cmd.Service, cmd.Username, cmd.Password)
		msg := fmt.Sprintf("Successfully saved credential for %s", cmd.Service)
		if err != nil {
			msg = fmt.Sprintf("Error saving username %s in service %s: %v\n", cmd.Username, cmd.Service, err)
		}
		return &model{
			vault:       vault,
			serviceName: cmd.Service,
			username:    cmd.Username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     msg,
			showSecret:  false,
			state:       &EndState{},
		}
	case cli.DeleteAction:
		err := vault.Delete(cmd.Service, cmd.Username)
		msg := fmt.Sprintf("Service %s deleted successfully.", cmd.Service)
		if err != nil {
			msg = fmt.Sprintf("Error deleting username %s in service %s: %v\n", cmd.Username, cmd.Service, err)
		}
		return &model{
			vault:       vault,
			serviceName: cmd.Service,
			username:    "", // user no longer exists after deletion
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     msg,
			showSecret:  false,
			state:       &EndState{},
		}
	case cli.RotateAction:
		err := vault.Replace(cmd.Service, cmd.Username, cmd.Password)
		msg := fmt.Sprintf("Service %s rotated successfully.", cmd.Service)
		if err != nil {
			msg = fmt.Sprintf("Error rotating username %s in service %s: %v\n", cmd.Username, cmd.Service, err)
		}
		return &model{
			vault:       vault,
			serviceName: cmd.Service,
			username:    cmd.Username,
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
			username:    cmd.Username,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  false,
			state:       &WelcomeState{},
		}
	}
}
