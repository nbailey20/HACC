package display

import (
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
