package engine

import (
	"fmt"

	"github.com/nbailey20/hacc/cli"
	"github.com/nbailey20/hacc/helpers"
	"github.com/nbailey20/hacc/vault"
)

type Executor struct {
	vault *vault.Vault
}

type Result struct {
	Action  string       `json:"action"`
	Success bool         `json:"success"`
	Error   string       `json:"error,omitempty"`
	Data    []Credential `json:"data,omitempty"`
}

type Credential struct {
	Service  string `json:"service"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewExecutor(vault *vault.Vault) *Executor {
	return &Executor{vault: vault}
}

func (e *Executor) Execute(cmd cli.CLICommand) Result {
	switch cmd.Action.(type) {
	case cli.AddAction:
		return e.add(cmd)
	case cli.DeleteAction:
		return e.delete(cmd)
	case cli.RotateAction:
		return e.rotate(cmd)
	case cli.BackupAction:
		return e.backup(cmd)
	case cli.SearchAction:
		return e.search(cmd)
	}
	return Result{
		Action:  "unknown",
		Success: false,
		Error:   fmt.Sprintf("unknown action: %T", cmd.Action),
	}
}

func (e *Executor) HasService(service string) bool {
	return e.vault.HasService(service)
}

func (e *Executor) HasUser(service, username string) bool {
	s, ok := e.vault.Services[service]
	if !ok {
		return false
	}
	return s.HasUser(username)
}

func (e *Executor) GetNumUsers(service string) int {
	s, ok := e.vault.Services[service]
	if !ok {
		return 0
	}
	return s.NumUsers()
}

func (e *Executor) GetUsersForService(service string) ([]string, error) {
	return e.vault.GetUsersForService(service)
}

func (e *Executor) GetNumServices() int {
	return len(e.vault.ListServices(""))
}

func (e *Executor) GetServices() []string {
	return e.vault.ListServices("")
}

func (e *Executor) GetServicesWithPrefix(prefix string) []string {
	return e.vault.ListServices(prefix)
}

func (e *Executor) GeneratePassword(digits, specials string, minLen, maxLen int) string {
	return helpers.GeneratePassword(digits, specials, minLen, maxLen)
}
