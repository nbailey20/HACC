package display

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nbailey20/hacc/hacc/helpers"
	"github.com/nbailey20/hacc/hacc/vault"
)

func addCredentialCmd(
	vault vault.Vault,
	service string,
	user string,
	pass string,
) tea.Cmd {
	return func() tea.Msg {
		if err := vault.Add(service, user, pass); err != nil {
			return AddFailedMsg{Error: err}
		}
		return nil
	}
}

func addMultiCredentialCmd(
	vault vault.Vault,
	file string,
) tea.Cmd {
	return func() tea.Msg {
		data, err := helpers.ReadJsonFile(file)
		if err != nil {
			return AddFailedMsg{Error: err}
		}
		creds := data["creds_list"]
		var addErrors []error
		for _, cred := range creds {
			currentError := vault.Add(
				cred["service"],
				cred["username"],
				cred["password"],
			)
			if currentError != nil {
				addErrors = append(addErrors, currentError)
			}
		}
		return AddFailedMsg{Error: errors.Join(addErrors...)}
	}
}

func deleteCredentialCmd(
	vault vault.Vault,
	service string,
	user string,
) tea.Cmd {
	return func() tea.Msg {
		if err := vault.Delete(service, user); err != nil {
			return DeleteFailedMsg{Error: err}
		}
		return nil
	}
}

func rotateCredentialCmd(
	vault vault.Vault,

	service string,
	user string,
	pass string,
) tea.Cmd {
	return func() tea.Msg {
		if err := vault.Replace(service, user, pass); err != nil {
			return RotateFailedMsg{Error: err}
		}
		return nil
	}
}

func generatePasswordCmd() tea.Cmd {
	return func() tea.Msg {
		return PasswordGeneratedMsg{Password: helpers.GeneratePassword()}
	}
}

func backupCredentialCmd(
	vault vault.Vault,
	fileName string,
	service string,
	user string,
) tea.Cmd {
	return func() tea.Msg {
		jsonData := map[string][]map[string]string{"creds_list": {}}
		if service != "" && user != "" {
			pass, err := vault.Get(service, user)
			if err != nil {
				return BackupErrorMsg{Error: err}
			}
			jsonData["creds_list"] = append(
				jsonData["creds_list"],
				map[string]string{
					"service":  service,
					"username": user,
					"password": pass,
				},
			)
			err = helpers.WriteJsonFile(fileName, jsonData)
			if err != nil {
				return BackupErrorMsg{Error: err}
			}
		}
		return nil
	}
}
