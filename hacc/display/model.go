package display

import (
	cli "github.com/nbailey20/hacc/hacc/cli"
	vault "github.com/nbailey20/hacc/hacc/vault"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	vault          vault.Vault
	state          State
	initialCmd     tea.Cmd
	action         cli.CLIAction
	serviceName    string
	username       string
	password       string
	digitsInPass   string
	specialsInPass string
	minPassLen     int
	maxPassLen     int
	page           int
	pageSize       int
	cursor         int
	showSecret     bool
	endSuccess     bool
	endMessage     string
	endError       error
}

func (m model) Init() tea.Cmd {
	return m.initialCmd
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

func initialState(cmd cli.CLICommand, vault vault.Vault) State {
	switch cmd.Action.(type) {
	case cli.SearchAction:
		return searchModelState(cmd.Service, cmd.Username, vault)
	case cli.AddAction:
		if cmd.Generate {
			return &ConfirmState{}
		}
		return &EndState{}
	case cli.DeleteAction:
		return &EndState{}
	case cli.RotateAction:
		return &EndState{}
	case cli.BackupAction:
		return &EndState{}
	default:
		return &WelcomeState{}
	}
}

func initialCmd(cmd cli.CLICommand, vault vault.Vault) tea.Cmd {
	switch cmd.Action.(type) {
	case cli.AddAction:
		return addCmd(vault, cmd)
	case cli.DeleteAction:
		return deleteCmd(vault, cmd)
	case cli.RotateAction:
		return rotateCmd(vault, cmd)
	case cli.BackupAction:
		return backupCmd(vault, cmd)
	default:
		return nil
	}
}

func NewModel(cmd cli.CLICommand, vault vault.Vault) *model {
	return &model{
		vault:          vault,
		state:          initialState(cmd, vault),
		serviceName:    cmd.Service,
		initialCmd:     initialCmd(cmd, vault),
		action:         cmd.Action,
		username:       cmd.Username,
		password:       cmd.Password,
		digitsInPass:   cmd.DigitsInPass,
		specialsInPass: cmd.SpecialsInPass,
		minPassLen:     cmd.MinLen,
		maxPassLen:     cmd.MaxLen,
		page:           0,
		pageSize:       pageSize,
		cursor:         0,
		endSuccess:     true,
	}
}
