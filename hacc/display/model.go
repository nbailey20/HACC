package display

import (
	tea "github.com/charmbracelet/bubbletea"

	cli "github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/engine"
)

type model struct {
	exec *engine.Executor
	cmd  cli.CLICommand

	state      State
	initialCmd tea.Cmd

	page     int
	pageSize int
	cursor   int
	showPass bool

	result engine.Result
}

func (m model) Init() tea.Cmd {
	return m.initialCmd
}

func searchModelState(service string, user string, exec *engine.Executor) State {
	// determines starting state for search commands
	// if we already have enough info from CLI/autocompletion,
	// jump straight to credential view
	if service == "" && user == "" {
		return &WelcomeState{}
	}
	if exec.HasService(service) && exec.HasUser(service, user) {
		return &CredentialState{}
	}
	if exec.HasService(service) {
		return &UsernameListState{}
	}
	return &ServiceListState{}
}

func initialState(cmd cli.CLICommand, exec *engine.Executor) State {
	switch cmd.Action.(type) {
	case cli.SearchAction:
		return searchModelState(cmd.Service, cmd.Username, exec)
	case cli.AddAction:
		if cmd.Generate {
			return &ConfirmState{}
		}
		return &EndState{}
	case cli.DeleteAction:
		return &EndState{}
	case cli.RotateAction:
		return &EndState{}
	case cli.BackupAction:
		return &EndState{}
	default:
		return &WelcomeState{}
	}
}

func initialCmd(cmd cli.CLICommand, exec *engine.Executor) tea.Cmd {
	switch cmd.Action.(type) {
	case cli.SearchAction:
		state := initialState(cmd, exec)
		if _, ok := state.(*CredentialState); ok {
			return execCmd(cmd, exec)
		}
		return nil
	case cli.AddAction:
		if cmd.Generate {
			return generateCmd(cmd, exec)
		}
		return execCmd(cmd, exec)
	case cli.DeleteAction:
		return execCmd(cmd, exec)
	case cli.RotateAction:
		return execCmd(cmd, exec)
	case cli.BackupAction:
		return execCmd(cmd, exec)
	default:
		return nil
	}
}

func NewModel(cmd cli.CLICommand, executor *engine.Executor) *model {
	st := initialState(cmd, executor)
	ic := initialCmd(cmd, executor)
	// If there's an initial command that will produce a result, enter a
	// pending state so the UI waits for the result rather than showing an
	// error/empty view immediately. Generate-only commands are local and
	// synchronous, so they should skip the loading state and render
	// immediately.
	if ic != nil {
		if _, ok := cmd.Action.(cli.AddAction); !ok || !cmd.Generate {
			st = &PendingState{Next: st}
		}
	}
	return &model{
		exec:       executor,
		cmd:        cmd,
		state:      st,
		initialCmd: ic,
		showPass:   false,
		page:       0,
		pageSize:   pageSize,
		cursor:     0,
	}
}
