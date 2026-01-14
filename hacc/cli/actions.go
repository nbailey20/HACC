package cli

import "github.com/spf13/cobra"

func registerActions(root *cobra.Command, cliCommand *CLICommand) {
	addSearchCommand(root, cliCommand)
	addAddCommand(root, cliCommand)
	addDeleteCommand(root, cliCommand)
	addRotateCommand(root, cliCommand)
}

func isVerb(s string) bool {
	switch s {
	case "search", "add", "a", "delete", "d", "rotate", "r", "backup", "b", "configure", "c":
		return true
	default:
		return false
	}
}

func addSearchCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		searchUsername string
	)

	searchCmd := &cobra.Command{
		Use:     "hacc SERVICE",
		Aliases: []string{"search"},
		Short:   "View credentials for a service",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCommand.Action = SearchAction{}
			cliCommand.Service = args[0]
			cliCommand.Username = searchUsername
			return nil
		},
	}

	searchCmd.Flags().StringVarP(&searchUsername, "username", "u", "", "Username")
	root.AddCommand(searchCmd)
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

func addRotateCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		rotateUsername string
		rotatePassword string
	)

	var rotateCmd = &cobra.Command{
		Use:     "rotate SERVICE",
		Aliases: []string{"d"},
		Short:   "rotate a credential for a service",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCommand.Action = RotateAction{}
			cliCommand.Service = args[0]
			cliCommand.Username = rotateUsername
			cliCommand.Password = rotatePassword
			return nil
		},
	}

	rotateCmd.Flags().StringVarP(&rotateUsername, "username", "u", "", "Username")
	rotateCmd.Flags().StringVarP(&rotatePassword, "password", "p", "", "Password")
	root.AddCommand(rotateCmd)
}

func addHelpCommand(root *cobra.Command, cliCommand *CLICommand) {
	var helpCmd = &cobra.Command{
		Use:     "help SERVICE",
		Aliases: []string{"d"},
		Short:   "help a credential for a service",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCommand.Action = HelpAction{}
			return nil
		},
	}
	root.AddCommand(helpCmd)
}
