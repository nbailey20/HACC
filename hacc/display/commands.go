package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/engine"
)

func execCmd(cmd cli.CLICommand, exec *engine.Executor) tea.Cmd {
	return func() tea.Msg {
		return resultMsg{exec.Execute(cmd)}
	}
}

func generateCmd(cmd cli.CLICommand, exec *engine.Executor) tea.Cmd {
	return func() tea.Msg {
		return PasswordGeneratedMsg{Password: exec.GeneratePassword(
			cmd.DigitsInPass,
			cmd.SpecialsInPass,
			cmd.MinLen,
			cmd.MaxLen,
		)}
	}
}

func copyPasswordCmd(password string) tea.Cmd {
	return func() tea.Msg {
		clipboard.Init()
		clipboard.Write(clipboard.FmtText, []byte(password))
		return PasswordCopiedMsg{}
	}
}
