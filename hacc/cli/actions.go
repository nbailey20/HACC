package cli

import "github.com/spf13/cobra"

func registerActions(root *cobra.Command, cliCommand *CLICommand) {
	addAddCommand(root, cliCommand)
	addDeleteCommand(root, cliCommand)
}

func isVerb(s string) bool {
	switch s {
	case "add", "a", "delete", "d", "rotate", "r", "backup", "b", "configure", "c":
		return true
	default:
		return false
	}
}

func addAddCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		addUsername string
		addPassword string
		addGenerate bool
	)

	addCmd := &cobra.Command{
		Use:     "add SERVICE",
		Aliases: []string{"a"},
		Short:   "Add a credential for a service",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCommand.Action = AddAction{}
			cliCommand.Service = args[0]
			cliCommand.Username = addUsername
			cliCommand.Password = addPassword
			cliCommand.Generate = addGenerate
			return nil
		},
	}

	addCmd.Flags().StringVarP(&addUsername, "username", "u", "", "Username")
	addCmd.Flags().StringVarP(&addPassword, "password", "p", "", "Password")
	addCmd.Flags().BoolVarP(&addGenerate, "generate", "g", false, "Generate password")
	root.AddCommand(addCmd)
}

func addDeleteCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		deleteUsername string
	)

	var deleteCmd = &cobra.Command{
		Use:     "delete SERVICE",
		Aliases: []string{"d"},
		Short:   "Delete a credential for a service",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCommand.Action = DeleteAction{}
			cliCommand.Service = args[0]
			cliCommand.Username = deleteUsername
			return nil
		},
	}

	deleteCmd.Flags().StringVarP(&deleteUsername, "username", "u", "", "Username")
	root.AddCommand(deleteCmd)
}
