package display

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/engine"
)

func Start(command cli.CLICommand, executor *engine.Executor) error {
	// If JSON output is requested, non-interactive mode is used.
	// If the command is not specific enough or it fails,
	//   an error field will be present in the response.
	if command.JSONOutput {
		result := executor.Execute(command)
		jsonBytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(jsonBytes))
		}
		return nil
	}

	// Default interactive mode
	numCreds := 0
	for _, service := range executor.GetServices() {
		numCreds += executor.GetNumUsers(service)
	}
	model := NewModel(command, executor)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
