package engine

import (
	"errors"
	"sync"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/helpers"
)

type backupResult struct {
	fileCred helpers.FileCred
	err      error
}

func (e *Executor) backup(cmd cli.CLICommand) Result {
	var resultsData []Credential
	var backupErrs []error
	var fileCreds []helpers.FileCred

	credsToBackup, err := e.getCredsForBackup(cmd)
	if err != nil {
		return Result{
			Action:  "backup",
			Success: false,
			Error:   err.Error(),
		}
	}
	totalCreds := 0
	for _, users := range credsToBackup {
		totalCreds += len(users)
	}

	// retrieve data for users in parallel
	results := make(chan backupResult, totalCreds)
	var wg sync.WaitGroup

	for service, users := range credsToBackup {
		for _, u := range users {
			wg.Add(1)
			go func(svc, user string) {
				defer wg.Done()
				fileCred, err := e.backupUser(svc, user)
				results <- backupResult{fileCred, err}
			}(service, u)
		}
	}
	// read the results of the channel
	for i := 0; i < totalCreds; i++ {
		result := <-results
		fileCreds = append(fileCreds, result.fileCred)
		resultsData = append(resultsData, Credential{Service: result.fileCred.Service, Username: result.fileCred.Username})
		if result.err != nil {
			backupErrs = append(backupErrs, result.err)
		}
	}

	// write data to creds file
	writeErr := helpers.WriteCredsFile(cmd.File, fileCreds)
	if writeErr != nil {
		return Result{
			Action:  "backup",
			Success: false,
			Error:   writeErr.Error(),
		}
	}
	if backupErrs != nil {
		return Result{
			Action:  "backup",
			Success: false,
			Error:   errors.Join(backupErrs...).Error(),
			Data:    resultsData,
		}
	}
	return Result{
		Action:  "backup",
		Success: true,
		Data:    resultsData,
	}
}

func (e *Executor) getCredsForBackup(cmd cli.CLICommand) (map[string][]string, error) {
	// determine users to backup
	credsToBackup := make(map[string][]string)
	var err error
	switch {
	case cmd.Username != "" && cmd.Service != "":
		credsToBackup[cmd.Service] = []string{cmd.Username}
		err = nil
	case cmd.Service != "":
		credsToBackup[cmd.Service], err = e.vault.GetUsersForService(cmd.Service)
	default:
		// backup entire Vault
		var allErrs []error
		for _, s := range e.vault.ListServices("") {
			credsToBackup[s], err = e.vault.GetUsersForService(s)
			if err != nil {
				allErrs = append(allErrs, err)
			}
		}
		err = errors.Join(allErrs...)
	}
	return credsToBackup, err
}

func (e *Executor) backupUser(service, user string) (helpers.FileCred, error) {
	pass, err := e.vault.Get(service, user)
	if err != nil {
		return helpers.FileCred{}, err
	}
	return helpers.FileCred{
		Service:  service,
		Username: user,
		Password: pass,
	}, nil
}
