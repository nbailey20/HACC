package main

import (
	"flag"
	"fmt"
)

type CLIInput struct {
	Action string
	Name   string
	Value  string
}

func parseInput(fs *flag.FlagSet, args []string) CLIInput {
	// Define flags
	action := fs.String("action", "get", "Action to perform: add, delete, replace, get (default)")
	name := fs.String("name", "", "Name of the service")
	value := fs.String("value", "", "Value for the service (for add/replace)")
	// Parse the flags
	fs.Parse(args)

	return CLIInput{
		Action: *action,
		Name:   *name,
		Value:  *value,
	}
}

func NewModel(cli CLIInput, vault Vault) *model {
	switch cli.Action {
	case "get":
		if cli.Name != "" {
			return &model{
				vault:       vault,
				serviceName: cli.Name,
				page:        0,
				pageSize:    pageSize,
				cursor:      0,
				message:     "",
				showSecret:  true,
				state:       &DetailState{},
			}
		}
		return &model{
			vault:       vault,
			serviceName: "",
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  true,
			state:       &WelcomeState{},
		}
	case "add":
		vault.Add(cli.Name, cli.Value)
		return &model{
			vault:       vault,
			serviceName: cli.Name,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     fmt.Sprintf("Successfully saved credential for %s", cli.Name),
			showSecret:  false,
			state:       &EndState{},
		}
	case "delete":
		err := vault.Delete(cli.Name)
		msg := fmt.Sprintf("Service %s deleted successfully.", cli.Name)
		if err != nil {
			msg = fmt.Sprintf("Error deleting service %s: %v\n", cli.Name, err)
		}
		return &model{
			vault:       vault,
			serviceName: "",
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     msg,
			showSecret:  false,
			state:       &EndState{},
		}
	default:
		return &model{
			vault:       vault,
			serviceName: "",
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  false,
			state:       &WelcomeState{},
		}
	}
}
