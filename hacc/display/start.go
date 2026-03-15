package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/engine"
)

func Start(command cli.CLICommand, executor *engine.Executor) error {
	numCreds := 0
	for _, service := range executor.GetServices() {
		numCreds += executor.GetNumUsers(service)
	}
	model := NewModel(command, executor)
	p := tea.NewProgram(model) //tea.WithAltScreen())
	_, err := p.Run()
	return err
}
