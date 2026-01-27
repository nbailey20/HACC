package display

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"

	"github.com/nbailey20/hacc/hacc/cli"
	"github.com/nbailey20/hacc/hacc/helpers"
	"github.com/nbailey20/hacc/hacc/vault"
)

// main add cmd invoked as initial cmd or state change
func addCmd(vault *vault.Vault, cmd cli.CLICommand) tea.Cmd {
	if cmd.File != "" {
		return addMultiCredentialCmd(vault, cmd.File)
	}
	if cmd.Generate {
		return generatePasswordCmd(
			cmd.DigitsInPass,
			cmd.SpecialsInPass,
			cmd.MinLen,
			cmd.MaxLen,
		)
	}
	return addCredentialCmd(vault, cmd.Service, cmd.Username, cmd.Password)
}

func addCredentialCmd(
	vault *vault.Vault,
	service string,
	user string,
	pass string,
) tea.Cmd {
	return func() tea.Msg {
		if err := vault.Add(service, user, pass); err != nil {
			return AddFailedMsg{
				Error: err,
				Display: fmt.Sprintf(
					"%s to add username %s for service %s.",
					failed,
					user,
					service,
				),
			}
		}
		return AddSuccessMsg{
			Display: fmt.Sprintf(
				"%s added username %s for service %s.",
				successfully,
				user,
				service,
			),
		}
	}
}

func addMultiCredentialCmd(vault *vault.Vault, file string) tea.Cmd {
	return func() tea.Msg {
		fileCreds, err := helpers.ReadCredsFile(file)
		if err != nil {
			return AddFailedMsg{
				Error: err,
				Display: fmt.Sprintf(
					"%s to read import file %s.",
					failed,
					file,
				),
			}
		}

		vaultResults := vault.AddMulti(fileCreds)
		var successResults []string
		var errorResults []error
		for _, result := range vaultResults {
			if result.Success == false {
				errorResults = append(errorResults, result.Err)
			} else {
				successResults = append(
					successResults,
					fmt.Sprintf(
						"%s added user %s for service %s.",
						successfully,
						result.Username,
						result.Service,
					),
				)
			}
		}

		if len(errorResults) > 0 {
			return AddFailedMsg{
				Error: errors.Join(errorResults...),
				Display: fmt.Sprintf(
					"%s\n%s to import one or more credentials from file %s",
					strings.Join(successResults, "\n"), // multiple creds can be added at once, show ones that succeeded
					failed,
					file,
				),
			}
		}
		return AddSuccessMsg{
			Display: strings.Join(successResults, "\n"),
		}
	}
}

func deleteCmd(vault *vault.Vault, cmd cli.CLICommand) tea.Cmd {
	return func() tea.Msg {
		if err := vault.Delete(cmd.Service, cmd.Username); err != nil {
			return DeleteErrorMsg{
				Error: err,
				Display: fmt.Sprintf(
					"%s to delete user %s from service %s.",
					failed,
					cmd.Username,
					cmd.Service,
				),
			}
		}
		return DeleteSuccessMsg{
			Display: fmt.Sprintf(
				"%s deleted user %s from service %s.",
				successfully,
				cmd.Username,
				cmd.Service,
			),
		}
	}
}

func rotateCmd(vault *vault.Vault, cmd cli.CLICommand) tea.Cmd {
	return func() tea.Msg {
		if err := vault.Replace(cmd.Service, cmd.Username, cmd.Password); err != nil {
			return RotateErrorMsg{
				Error: err,
				Display: fmt.Sprintf(
					"%s to rotate user %s in service %s.",
					failed,
					cmd.Username,
					cmd.Service,
				),
			}
		}
		return RotateSuccessMsg{
			Display: fmt.Sprintf(
				"%s rotated user %s in service %s.",
				successfully,
				cmd.Username,
				cmd.Service,
			),
		}
	}
}

type backupResult struct {
	fileCred helpers.FileCred
	result   string
	err      error
}

func backupCmd(vault *vault.Vault, cmd cli.CLICommand) tea.Cmd {
	return func() tea.Msg {
		var backupResults []string
		var backupErrs []error
		var fileCreds []helpers.FileCred

		credsToBackup, err := getCredsForBackup(vault, cmd)
		if err != nil {
			return BackupErrorMsg{
				Error: err,
				Display: fmt.Sprintf(
					"%s retrieving credentials for backup.",
					failed,
				),
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
					fileCred, result, err := backupUserCmd(
						vault,
						svc,
						user,
					)
					results <- backupResult{fileCred, result, err}
				}(service, u)
			}
		}
		// read the results of the channel
		for i := 0; i < totalCreds; i++ {
			result := <-results
			fileCreds = append(fileCreds, result.fileCred)
			backupResults = append(backupResults, result.result)
			if result.err != nil {
				backupErrs = append(backupErrs, result.err)
			}
		}

		// write data to creds file
		writeErr := helpers.WriteCredsFile(cmd.File, fileCreds)
		if writeErr != nil {
			return BackupErrorMsg{
				Error: writeErr,
				Display: fmt.Sprintf(
					"%s to write credentials to file %s.",
					failed,
					cmd.File,
				),
			}
		}
		if backupErrs != nil {
			return BackupErrorMsg{
				Error: errors.Join(backupErrs...),
				Display: fmt.Sprintf(
					"%s\n\n%s to backup one or more credentials to file %s.",
					strings.Join(backupResults, "\n"),
					failed,
					cmd.File,
				),
			}
		}
		return BackupSuccessMsg{
			Display: fmt.Sprintf(
				"%s\n\n%s completed backup to file %s.",
				strings.Join(backupResults, "\n"),
				successfully,
				cmd.File,
			),
		}
	}
}

func getCredsForBackup(vault *vault.Vault, cmd cli.CLICommand) (map[string][]string, error) {
	// determine users to backup
	credsToBackup := make(map[string][]string)
	var err error
	switch {
	case cmd.Username != "" && cmd.Service != "":
		credsToBackup[cmd.Service] = []string{cmd.Username}
		err = nil
	case cmd.Service != "":
		credsToBackup[cmd.Service], err = vault.GetUsersForService(cmd.Service)
	default:
		// backup entire Vault
		var allErrs []error
		for _, s := range vault.ListServices("") {
			credsToBackup[s], err = vault.GetUsersForService(s)
			if err != nil {
				allErrs = append(allErrs, err)
			}
		}
		err = errors.Join(allErrs...)
	}
	return credsToBackup, err
}

func backupUserCmd(
	vault *vault.Vault,
	service string,
	user string,
) (helpers.FileCred, string, error) {
	successDisplay := fmt.Sprintf(
		"%s backed up credential %s/%s.",
		successfully,
		service,
		user,
	)
	errorDisplay := fmt.Sprintf(
		"%s to back up credential %s/%s.",
		failed,
		service,
		user,
	)
	pass, err := vault.Get(service, user)
	if err != nil {
		return helpers.FileCred{}, errorDisplay, err
	}
	return helpers.FileCred{
		Service:  service,
		Username: user,
		Password: pass,
	}, successDisplay, nil
}

func generatePasswordCmd(
	digits string,
	specials string,
	minLen int,
	maxLen int,
) tea.Cmd {
	return func() tea.Msg {
		return PasswordGeneratedMsg{
			Password: helpers.GeneratePassword(
				digits,
				specials,
				minLen,
				maxLen,
			),
		}
	}
}

func loadPasswordCmd(vault *vault.Vault, service string, user string) tea.Cmd {
	return func() tea.Msg {
		pass, err := vault.Get(service, user)
		if err != nil {
			return PasswordLoadedMsg{
				Password: fmt.Sprintf("Could not load password: %v", err),
			}
		}
		return PasswordLoadedMsg{
			Password: pass,
		}
	}
}

func copyPasswordCmd(password string) tea.Cmd {
	return func() tea.Msg {
		clipboard.Init()
		clipboard.Write(clipboard.FmtText, []byte(password))
		return PasswordCopiedMsg{}
	}
}
