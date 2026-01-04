package main

import (
	"flag"
	"fmt"
	"reflect"
)

type FlagDef struct {
	Short string
	Long  string
	Help  string
	Kind  string // "bool" or "string"
}

type ActionDef struct {
	Name string
	Flag FlagDef
}

var Actions = []ActionDef{
	{
		Name: "install",
		Flag: FlagDef{Short: "i", Long: "install", Help: "Create new authentication credential Vault", Kind: "bool"},
	},
	{
		Name: "eradicate",
		Flag: FlagDef{Short: "e", Long: "eradicate", Help: "Delete entire Vault - cannot be undone", Kind: "bool"},
	},
	{
		Name: "add",
		Flag: FlagDef{Short: "a", Long: "add", Help: "Add new set of credentials to Vault", Kind: "bool"},
	},
	{
		Name: "delete",
		Flag: FlagDef{Short: "d", Long: "delete", Help: "Delete credentials in Vault", Kind: "bool"},
	},
	{
		Name: "rotate",
		Flag: FlagDef{Short: "r", Long: "rotate", Help: "Rotate existing password in Vault", Kind: "bool"},
	},
	{
		Name: "backup",
		Flag: FlagDef{Short: "b", Long: "backup", Help: "Backup entire Vault", Kind: "bool"},
	},
	{
		Name: "configure",
		Flag: FlagDef{Short: "c", Long: "configure", Help: "View or modify client configuration", Kind: "bool"},
	},
	{
		Name: "upgrade",
		Flag: FlagDef{Short: "", Long: "upgrade", Help: "Upgrade client software", Kind: "bool"},
	},
}

var SubArgs = []FlagDef{
	{Short: "", Long: "service", Help: "Service to perform action on", Kind: "string"},
	{Short: "u", Long: "username", Help: "Username to perform action on", Kind: "string"},
	{Short: "p", Long: "password", Help: "Password for new credentials, used with add action", Kind: "string"},
	{Short: "g", Long: "generate", Help: "Generate random password for add operation", Kind: "bool"},
	{Short: "f", Long: "file", Help: "File name for import / export operations", Kind: "string"},
	{Short: "w", Long: "wipe", Help: "Wipe all existing credentials during Vault eradication", Kind: "bool"},
	{Short: "", Long: "export", Help: "Export existing client configuration as encrypted file", Kind: "bool"},
	{Short: "", Long: "set", Help: "Set client configuration parameter", Kind: "string"},
	{Short: "", Long: "show", Help: "Show client configuration parameter", Kind: "string"},
}

var AllowedInputCombinations = map[string][][]string{
	"install": {
		// {},
		// {"file"},
	},
	"eradicate": {
		// {},
		// {"wipe"},
		// {"upgrade"},
	},
	"add": {
		// {},
		// {"service"},
		// {"username"},
		// {"password"},
		// {"generate"},
		// {"file"},
		// {"service", "username"},
		// {"service", "password"},
		// {"service", "generate"},
		// {"username", "password"},
		// {"username", "generate"},
		{"service", "username", "password"},
		// {"service", "username", "generate"},
	},
	"delete": {
		// {},
		// {"service"},
		// {"username"},
		{"service", "username"},
	},
	"rotate": {
		//{},
		//{"service"},
		//{"username"},
		// {"password"},
		// {"generate"},
		// {"service", "username"},
		// {"service", "password"},
		// {"service", "generate"},
		// {"username", "password"},
		// {"username", "generate"},
		{"service", "username", "password"},
		// {"service", "username", "generate"},
	},

	"search": {
		{},
		// {"username"},
		{"service"},
		{"service", "username"},
	},
	"backup": {
		//{},
		//{"file"},
	},
	"configure": {
		// {"show"},
		// {"set"},
		//{"export"},
		// {"set", "file", "password"},
		// {"export", "file"},
	},
	"upgrade": {
		// {},
	},
}

type CLIInput struct {
	install   bool
	eradicate bool
	add       bool
	delete    bool
	rotate    bool
	backup    bool
	configure bool
	upgrade   bool
	service   string
	username  string
	password  string
	generate  bool
	file      string
	wipe      bool
	export    bool
	set       string
	show      string
}

type RawCLICommand struct {
	action   []string
	service  string
	username string
	password string
	generate bool
	file     string
	wipe     bool
	export   bool
	set      string
	show     string
}

type CLICommand struct {
	action   string
	service  string
	username string
	password string
	generate bool
	file     string
	wipe     bool
	export   bool
	set      string
	show     string
}

func parseInput(fs *flag.FlagSet, args []string) (CLIInput, error) {
	var cli CLIInput

	boolFlags := map[string]*bool{
		"install":   &cli.install,
		"eradicate": &cli.eradicate,
		"add":       &cli.add,
		"delete":    &cli.delete,
		"rotate":    &cli.rotate,
		"backup":    &cli.backup,
		"configure": &cli.configure,
		"upgrade":   &cli.upgrade,
	}

	stringFlags := map[string]*string{
		"service":  &cli.service,
		"username": &cli.username,
		"password": &cli.password,
		"file":     &cli.file,
		"set":      &cli.set,
		"show":     &cli.show,
	}

	// Register long flags for actions dynamically
	for _, action := range Actions {
		if action.Flag.Kind == "bool" {
			if ptr := boolFlags[action.Name]; ptr != nil {
				fs.BoolVar(ptr, action.Flag.Long, false, action.Flag.Help)
			}
		}
	}
	// Register short flags for actions
	for _, action := range Actions {
		if action.Flag.Kind == "bool" && action.Flag.Short != "" {
			if ptr := boolFlags[action.Name]; ptr != nil {
				fs.BoolVar(ptr, action.Flag.Short, false, action.Flag.Help)
			}
		}
	}

	// Register long flags for sub-arguments
	for _, subArg := range SubArgs {
		switch subArg.Kind {
		case "bool":
			if ptr := boolFlags[subArg.Long]; ptr != nil {
				fs.BoolVar(ptr, subArg.Long, false, subArg.Help)
			}
		case "string":
			if ptr := stringFlags[subArg.Long]; ptr != nil {
				fs.StringVar(ptr, subArg.Long, "", subArg.Help)
			}
		}
	}
	// Register available short flags for sub-arguments
	for _, subArg := range SubArgs {
		if subArg.Short != "" {
			switch subArg.Kind {
			case "bool":
				if ptr := boolFlags[subArg.Long]; ptr != nil {
					fs.BoolVar(ptr, subArg.Short, false, subArg.Help)
				}
			case "string":
				if ptr := stringFlags[subArg.Long]; ptr != nil {
					fs.StringVar(ptr, subArg.Short, "", subArg.Help)
				}
			}
		}
	}

	// Parse the flags
	if err := fs.Parse(args); err != nil {
		return cli, err
	}

	// After parsing, fs.Args() contains positional arguments
	if len(fs.Args()) > 0 {
		cli.service = fs.Args()[0] // capture first positional arg as service
	}
	return cli, nil
}

func normalizeInput(input CLIInput) RawCLICommand {
	var action []string
	if input.install {
		action = append(action, "install")
	}
	if input.eradicate {
		action = append(action, "eradicate")
	}
	if input.add {
		action = append(action, "add")
	}
	if input.delete {
		action = append(action, "delete")
	}
	if input.rotate {
		action = append(action, "rotate")
	}
	if input.backup {
		action = append(action, "backup")
	}
	if input.configure {
		action = append(action, "configure")
	}
	if input.upgrade {
		action = append(action, "upgrade")
	}

	if len(action) == 0 {
		action = append(action, "search") // default action
	}

	return RawCLICommand{
		action:   action,
		service:  input.service,
		username: input.username,
		password: input.password,
		generate: input.generate,
		file:     input.file,
		wipe:     input.wipe,
		export:   input.export,
		set:      input.set,
		show:     input.show,
	}
}

func validateCommand(input RawCLICommand) (CLICommand, error) {
	// ensure only one action is specified
	numActionFlags := len(input.action)
	if numActionFlags > 1 {
		return CLICommand{}, fmt.Errorf("Multiple actions specified, only one action allowed")
	}

	// ensure a valid action is specified
	actionName := input.action[0]
	allowedCombinations, exists := AllowedInputCombinations[actionName]
	if !exists {
		return CLICommand{}, fmt.Errorf("Unknown action: %s", actionName)
	}

	// determine which subargs are provided
	var providedSubargs []string
	for _, subArg := range SubArgs {
		switch subArg.Kind {
		case "bool":
			v := reflect.ValueOf(input).FieldByName(subArg.Long)
			if v.IsValid() && v.Bool() {
				providedSubargs = append(providedSubargs, subArg.Long)
			}
		case "string":
			v := reflect.ValueOf(input).FieldByName(subArg.Long)
			if v.IsValid() && v.String() != "" {
				providedSubargs = append(providedSubargs, subArg.Long)
			}
		}
	}

	// Check if all provided subargs are in an allowed combination for the action
	validCombination := false
	for _, combo := range allowedCombinations {
		if equalStringSets(combo, providedSubargs) {
			validCombination = true
			break
		}
	}
	if !validCombination {
		return CLICommand{}, fmt.Errorf("Invalid combination of flags")
	}
	return CLICommand{
		action:   actionName,
		service:  input.service,
		username: input.username,
		password: input.password,
		generate: input.generate,
		file:     input.file,
		wipe:     input.wipe,
		export:   input.export,
		set:      input.set,
		show:     input.show,
	}, nil
}

// Helper function to compare two string slices as sets (ignoring order)
func equalStringSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	seen := make(map[string]int, len(a))
	for _, v := range a {
		seen[v]++
	}

	for _, v := range b {
		if seen[v] == 0 {
			return false
		}
		seen[v]--
	}

	return true
}

func autoCompleteCommand(command CLICommand, vault *Vault) CLICommand {
	// Auto-complete missing fields in command where possible using vault data
	command = autoCompleteService(command, vault)
	command = autoCompleteUsername(command, vault)
	return command
}
func autoCompleteService(command CLICommand, vault *Vault) CLICommand {
	// If only one service exists in vault for a given prefix, auto-complete it
	if len(vault.ListServices(command.service)) == 1 {
		command.service = vault.ListServices(command.service)[0]
	}
	return command
}

func autoCompleteUsername(command CLICommand, vault *Vault) CLICommand {
	// If only one user prefix exists for a valid service, auto-complete it
	if vault.HasService(command.service) {
		if len(vault.services[command.service].GetUsers(command.username)) == 1 {
			command.username = vault.services[command.service].GetUsers(command.username)[0]
		}
	}
	return command
}
