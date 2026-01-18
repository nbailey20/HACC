package cli

import (
	"fmt"
	"reflect"
	"unicode"
)

type FlagDef struct {
	Short string
	Long  string
	Help  string
	Kind  string // "bool" or "string"
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

var AllowedInputCombinations = map[ActionKind][][]string{
	ActionInstall: {
		// {},
		// {"file"},
	},
	ActionEradicate: {
		// {},
		// {"wipe"},
		// {"upgrade"},
	},
	ActionAdd: {
		{"file"},
		{"service", "username", "password"},
		{"service", "username", "generate"},
	},
	ActionDelete: {
		{"service", "username"},
	},
	ActionRotate: {
		{"service", "username", "password"},
		{"service", "username", "generate"},
	},
	ActionSearch: {
		{},
		{"service"},
		{"service", "username"},
	},
	ActionBackup: {
		{"file"},
		{"file", "service"},
		{"file", "service", "username"},
	},
	ActionConfigure: {
		// {"show"},
		// {"set"},
		//{"export"},
		// {"set", "file", "password"},
		// {"export", "file"},
	},
	ActionUpgrade: {
		// {},
	},
}

func uppercaseFirst(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func ValidateCommand(command CLICommand) error {
	// ensure a valid action is specified
	if command.Action == nil {
		return fmt.Errorf("No action")
	}
	actionType := command.Action.Kind()
	allowedCombinations, exists := AllowedInputCombinations[actionType]
	if !exists {
		return fmt.Errorf("Unknown action")
	}

	// determine which subargs are provided
	var providedSubargs []string
	for _, subArg := range SubArgs {
		switch subArg.Kind {
		case "bool":
			v := reflect.ValueOf(command).FieldByName(uppercaseFirst(subArg.Long))
			if v.IsValid() && v.Bool() {
				providedSubargs = append(providedSubargs, subArg.Long)
			}
		case "string":
			v := reflect.ValueOf(command).FieldByName(uppercaseFirst(subArg.Long))
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
		return fmt.Errorf("Invalid combination of flags: %v", providedSubargs)
	}
	return nil
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
