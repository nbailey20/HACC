package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/vault"
)

func Start(command cli.CLICommand, vault *vault.Vault) error {
	numCreds := 0
	for _, service := range vault.Services {
		numCreds += service.NumUsers()
	}

	model := NewModel(command, vault)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
