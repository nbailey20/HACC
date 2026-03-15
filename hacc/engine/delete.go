package engine

import (
	"errors"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/helpers"
)

func (e *Executor) delete(cmd cli.CLICommand) Result {
	if cmd.File != "" {
		return e.deleteMultiCredential(cmd.File)
	}
	return e.deleteCredential(cmd)
}

// single delete cmd wrapper used when no file provided
func (e *Executor) deleteCredential(cmd cli.CLICommand) Result {
	err := e.vault.Delete(cmd.Service, cmd.Username)
	result := Result{
		Action:  "delete",
		Success: err == nil,
		Data:    []Credential{{Service: cmd.Service, Username: cmd.Username}},
	}
	if err != nil {
		result.Error = err.Error()
	}
	return result
}

// delete multiple credentials in parallel from a file
func (e *Executor) deleteMultiCredential(file string) Result {
	fileCreds, err := helpers.ReadCredsFile(file)
	if err != nil {
		return Result{
			Action:  "delete",
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		}
	}

	results := e.vault.DeleteMulti(fileCreds)
	var successfulCreds []Credential
	var errorResults []error
	for _, result := range results {
		if result.Success == false {
			errorResults = append(errorResults, result.Err)
		} else {
			successfulCreds = append(
				successfulCreds,
				Credential{Service: result.Service, Username: result.Username},
			)
		}
	}

	result := Result{
		Action:  "delete",
		Success: len(errorResults) == 0,
		Data:    successfulCreds,
	}
	if len(errorResults) > 0 {
		result.Error = errors.Join(errorResults...).Error()
	}
	return result
}
