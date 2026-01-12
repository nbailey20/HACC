package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func preprocessArgs(args []string) []string {
	// hacc add gmail -u test -g is supported by default
	// also support: hacc gmail add -u test -g
	// don't forget shorthand: hacc gmail a -u test -g
	if len(args) > 1 && isVerb(args[1]) {
		args[0], args[1] = args[1], args[0]
	}
	return args
}

func NewRootCommand(cliCommand *CLICommand) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hacc",
		Short: "hacc is a credential manager backed by AWS SSM",
		Args:  cobra.ArbitraryArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cliCommand.Action = SearchAction{}
				return nil
			}
			if len(args) == 1 {
				cliCommand.Action = SearchAction{}
				cliCommand.Service = args[0]
				return nil
			}
			return fmt.Errorf("Unrecognized action in %v", args)
		},
	}

	registerActions(rootCmd, cliCommand)
	return rootCmd
}

// Parse runs the cobra command tree and returns the populated CLICommand.
func Parse(args []string, quiet bool) (CLICommand, error) {
	// ensure args are in the order expected by Cobra
	args = preprocessArgs(args)

	var cliCommand CLICommand
	rootCmd := NewRootCommand(&cliCommand)
	rootCmd.SetArgs(args)

	// don't print usage instructions during tests
	if quiet {
		rootCmd.SilenceUsage = true
		rootCmd.SilenceErrors = true
	}

	if err := rootCmd.Execute(); err != nil {
		return CLICommand{}, err
	}

	return cliCommand, nil
}
