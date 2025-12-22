package main

import (
	"flag"
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

// func main() {
// 	fs := flag.NewFlagSet("main", flag.ExitOnError)
// 	input := parseInput(fs, os.Args[1:]) // Pass the command-line arguments
// 	fmt.Printf("Action: %s, Name: %s, Value: %s\n", input.Action, input.Name, input.Value)
// }
