package display

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"

	"github.com/nbailey20/hacc/hacc/cli"
	"github.com/nbailey20/hacc/hacc/helpers"
	"github.com/nbailey20/hacc/hacc/vault"
)

func addCmd(vault vault.Vault, cmd cli.CLICommand) tea.Cmd {
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
	vault vault.Vault,
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

func addMultiCredentialCmd(vault vault.Vault, file string) tea.Cmd {
	return func() tea.Msg {
		data, err := helpers.ReadJsonFile(file)
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
		creds := data["creds_list"]
		var addErrors []error
		var addSuccesses []string
		for _, cred := range creds {
			currentError := vault.Add(
				cred["service"],
				cred["username"],
				cred["password"],
			)
			if currentError != nil {
				addErrors = append(addErrors, currentError)
			} else {
				addSuccesses = append(
					addSuccesses,
					fmt.Sprintf(
						"%s added user %s for service %s.",
						successfully,
						cred["username"],
						cred["service"],
					),
				)
			}
		}
		if len(addErrors) > 0 {
			return AddFailedMsg{
				Error: errors.Join(addErrors...),
				Display: fmt.Sprintf(
					"%s\n%s to import one or more credentials from file %s",
					strings.Join(addSuccesses, "\n"), // multiple creds can be added at once, show ones that succeeded
					failed,
					file,
				),
			}
		}
		return AddSuccessMsg{
			Display: strings.Join(addSuccesses, "\n"),
		}
	}
}

func deleteCmd(vault vault.Vault, cmd cli.CLICommand) tea.Cmd {
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

func rotateCmd(vault vault.Vault, cmd cli.CLICommand) tea.Cmd {
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

func backupCmd(vault vault.Vault, cmd cli.CLICommand) tea.Cmd {
	return func() tea.Msg {
		var jsonData = map[string][]map[string]string{"creds_list": {}}
		var display string
		var backupErr error
		switch {
		case cmd.Username != "" && cmd.Service != "":
			var serviceJson map[string]string
			serviceJson, display, backupErr = backupUserCmd(
				vault,
				cmd.Service,
				cmd.Username,
			)
			jsonData["creds_list"] = []map[string]string{serviceJson}
		case cmd.Service != "":
			var userDisplays []string
			var userErrors []error
			jsonData["creds_list"], userDisplays, userErrors = backupServiceCmd(
				vault,
				cmd.Service,
			)
			display = strings.Join(userDisplays, "\n")
			backupErr = errors.Join(userErrors...)
		default:
			var serviceDisplays []string
			var serviceErrors []error
			jsonData["creds_list"], serviceDisplays, serviceErrors = backupVaultCmd(
				vault,
			)
			display = strings.Join(serviceDisplays, "\n")
			backupErr = errors.Join(serviceErrors...)
		}

		writeErr := helpers.WriteJsonFile(cmd.File, jsonData)
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
		if backupErr != nil {
			return BackupErrorMsg{
				Error: backupErr,
				Display: fmt.Sprintf(
					"%s\n\n%s to backup one or more credentials to file %s.",
					display,
					failed,
					cmd.File,
				),
			}
		}
		return BackupSuccessMsg{
			Display: fmt.Sprintf(
				"%s\n\n%s completed backup to file %s.",
				display,
				successfully,
				cmd.File,
			),
		}

	}
}

func backupVaultCmd(vault vault.Vault) ([]map[string]string, []string, []error) {
	services := vault.ListServices("")
	// aggregate json/errors/display results from lower cmd
	var backupJson []map[string]string
	var backupDisplays []string
	var backupErrors []error
	for _, service := range services {
		serviceJson, serviceDisplays, serviceErrs := backupServiceCmd(
			vault,
			service,
		)
		if serviceErrs != nil {
			backupErrors = append(backupErrors, serviceErrs...)
		} else {
			backupJson = append(backupJson, serviceJson...)
		}
		backupDisplays = append(backupDisplays, serviceDisplays...)
	}
	return backupJson, backupDisplays, backupErrors
}

func backupServiceCmd(
	vault vault.Vault,
	service string,
) ([]map[string]string, []string, []error) {
	users := vault.Services[service].GetUsers("")
	// aggregate json/errors/display results from lower cmd
	var backupJson []map[string]string
	var backupErrors []error
	var backupDisplays []string
	for _, user := range users {
		json, display, err := backupUserCmd(
			vault,
			service,
			user,
		)
		if err != nil {
			backupErrors = append(backupErrors, err)
		} else {
			backupJson = append(backupJson, json)
		}
		backupDisplays = append(backupDisplays, display)
	}
	return backupJson, backupDisplays, backupErrors
}

func backupUserCmd(
	vault vault.Vault,
	service string,
	user string,
) (map[string]string, string, error) {
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
		return nil, errorDisplay, err
	}
	return map[string]string{
		"service":  service,
		"username": user,
		"password": pass,
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

func loadPasswordCmd(vault vault.Vault, service string, user string) tea.Cmd {
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
