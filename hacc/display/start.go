package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nbailey20/hacc/hacc/cli"
	"github.com/nbailey20/hacc/hacc/vault"
)

func Start(command cli.CLICommand, vault *vault.Vault) error {
	model := NewModel(command, vault)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
