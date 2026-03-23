package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/engine"
)

func execCmd(cmd cli.CLICommand, exec *engine.Executor) tea.Cmd {
	return func() tea.Msg {
		result := exec.Execute(cmd)
		// if result is a single password, copy to clipboard
		if result.Action == "search" && len(result.Data) == 1 {
			copyPasswordCmd(result.Data[0].Password)
		}
		return resultMsg{result}
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

func copyPasswordCmd(password string) {
	clipboard.Init()
	clipboard.Write(clipboard.FmtText, []byte(password))
}
