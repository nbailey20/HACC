package engine

import (
	"github.com/nbailey20/hacc/cli"
)

func (e *Executor) rotate(cmd cli.CLICommand) Result {
	err := e.vault.Replace(cmd.Service, cmd.Username, cmd.Password)
	result := Result{
		Action:  "rotate",
		Success: err == nil,
		Data:    []Credential{{Service: cmd.Service, Username: cmd.Username}},
	}
	if err != nil {
		result.Error = err.Error()
	}
	return result
}
