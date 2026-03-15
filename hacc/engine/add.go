package engine

import (
	"errors"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/helpers"
)

func (e *Executor) add(cmd cli.CLICommand) Result {
	if cmd.File != "" {
		return e.addMultiCredential(cmd.File)
	}
	if cmd.Generate {
		return Result{
			Action:  "generate",
			Success: true,
			Data: []Credential{{Password: helpers.GeneratePassword(
				cmd.DigitsInPass,
				cmd.SpecialsInPass,
				cmd.MinLen,
				cmd.MaxLen),
			}},
		}
	}
	return e.addCredential(cmd.Service, cmd.Username, cmd.Password)
}

func (e *Executor) addCredential(service, user, pass string) Result {
	err := e.vault.Add(service, user, pass)
	result := Result{
		Action:  "add",
		Success: err == nil,
		Data:    []Credential{{Service: service, Username: user}},
	}
	if err != nil {
		result.Error = err.Error()
	}
	return result
}

func (e *Executor) addMultiCredential(file string) Result {
	fileCreds, err := helpers.ReadCredsFile(file)
	if err != nil {
		return Result{
			Action:  "add",
			Success: false,
			Error:   err.Error(),
			Data:    nil,
		}
	}

	vaultResults := e.vault.AddMulti(fileCreds)
	var successfulCreds []Credential
	var errorResults []error
	for _, result := range vaultResults {
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
		Action:  "add",
		Success: len(errorResults) == 0,
		Data:    successfulCreds,
	}
	if len(errorResults) > 0 {
		result.Error = errors.Join(errorResults...).Error()
	}
	return result
}
