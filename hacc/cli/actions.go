package cli

import "github.com/spf13/cobra"

func registerActions(root *cobra.Command, cliCommand *CLICommand) {
	addSearchCommand(root, cliCommand)
	addAddCommand(root, cliCommand)
	addDeleteCommand(root, cliCommand)
	addRotateCommand(root, cliCommand)
	addBackupCommand(root, cliCommand)
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
		addFile     string
	)

	addCmd := &cobra.Command{
		Use:     "add SERVICE",
		Aliases: []string{"a"},
		Short:   "Add a credential for a service",
		Args: func(cmd *cobra.Command, args []string) error {
			// if using file flag, no positional service arg expected
			if addFile != "" {
				return nil
			}
			// otherwise expect the service positional arg
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				cliCommand.Service = args[0]
			}
			cliCommand.Action = AddAction{}
			cliCommand.Username = addUsername
			cliCommand.Password = addPassword
			cliCommand.Generate = addGenerate
			cliCommand.File = addFile
			return nil
		},
	}

	addCmd.Flags().StringVarP(&addUsername, "username", "u", "", "Username")
	addCmd.Flags().StringVarP(&addPassword, "password", "p", "", "Password")
	addCmd.Flags().BoolVarP(&addGenerate, "generate", "g", false, "Generate password")
	addCmd.Flags().StringVarP(&addFile, "file", "f", "", "Backup file name to import")
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
		Aliases: []string{"r"},
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

func addBackupCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		backupFile string
		backupUser string
	)

	var backupCmd = &cobra.Command{
		Use:     "backup [SERVICE]",
		Aliases: []string{"b"},
		Short:   "backup credential for all or one service",
		Args: func(cmd *cobra.Command, args []string) error {
			// if no positional arg, backup everything
			if len(args) == 0 {
				return nil
			}
			// otherwise only backup the specific service/user
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCommand.Action = BackupAction{}
			if len(args) > 0 {
				cliCommand.Service = args[0]
			}
			cliCommand.Username = backupUser
			cliCommand.File = backupFile
			return nil
		},
	}

	backupCmd.Flags().StringVarP(&backupFile, "file", "f", "", "File name where backup should be created")
	backupCmd.Flags().StringVarP(&backupUser, "username", "u", "", "Username")
	root.AddCommand(backupCmd)
}
