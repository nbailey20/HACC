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
		addUsername    string
		addPassword    string
		addFile        string
		addGenerate    bool
		digitsInPass   string
		specialsInPass string
		minLen         int
		maxLen         int
		jsonOutput     bool
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
			cliCommand.File = addFile
			cliCommand.Generate = addGenerate
			cliCommand.SpecialsInPass = specialsInPass
			cliCommand.DigitsInPass = digitsInPass
			cliCommand.MinLen = minLen
			cliCommand.MaxLen = maxLen
			return nil
		},
	}

	addCmd.Flags().StringVarP(&addUsername, "username", "u", "", "Username")
	addCmd.Flags().StringVarP(&addPassword, "password", "p", "", "Password")
	addCmd.Flags().StringVarP(&addFile, "file", "f", "", "Backup file name for bulk credential addition")
	addCmd.Flags().BoolVarP(&addGenerate, "generate", "g", false, "Generate password")
	addCmd.Flags().StringVarP(&specialsInPass, "specials", "", "any", "Digit requirement when generating a password: any|required|forbidden")
	addCmd.Flags().StringVarP(&digitsInPass, "digits", "", "any", "Special character requirement when generating a password: any|required|forbidden")
	addCmd.Flags().IntVarP(&minLen, "min-len", "", 0, "Minimum length of password to generate")
	addCmd.Flags().IntVarP(&maxLen, "max-len", "", 0, "Maximum length of password to generate")
	addCmd.Flags().BoolVarP(&jsonOutput, "json-output", "", false, "Whether output should be json or interactive (default interactive)")

	// Flag constraints
	addCmd.MarkFlagsMutuallyExclusive("password", "generate")
	addCmd.MarkFlagsOneRequired("username", "file")
	addCmd.MarkFlagsOneRequired("password", "generate", "file")
	for _, a := range []string{
		"username",
		"password",
		"generate",
		"specials",
		"digits",
		"min-len",
		"max-len",
	} {
		addCmd.MarkFlagsMutuallyExclusive(a, "file")
	}
	root.AddCommand(addCmd)
}

func addDeleteCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		deleteUsername string
		deleteFile     string
		jsonOutput     bool
	)

	var deleteCmd = &cobra.Command{
		Use:     "delete SERVICE",
		Aliases: []string{"d"},
		Short:   "Delete a credential for a service",
		Args: func(cmd *cobra.Command, args []string) error {
			// if using file flag, no positional service arg expected
			if deleteFile != "" {
				return nil
			}
			// otherwise expect the service positional arg
			return cobra.ExactArgs(1)(cmd, args)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				cliCommand.Service = args[0]
			}
			cliCommand.Action = DeleteAction{}
			cliCommand.Username = deleteUsername
			cliCommand.File = deleteFile
			return nil
		},
	}

	deleteCmd.Flags().StringVarP(&deleteUsername, "username", "u", "", "Username")
	deleteCmd.Flags().StringVarP(&deleteFile, "file", "f", "", "Backup file name for bulk credential deletion")
	deleteCmd.Flags().BoolVarP(&jsonOutput, "json-output", "", false, "Whether output should be json or interactive (default interactive)")
	deleteCmd.MarkFlagsMutuallyExclusive("username", "file")
	deleteCmd.MarkFlagsOneRequired("username", "file")
	root.AddCommand(deleteCmd)
}

func addRotateCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		rotateUsername string
		rotatePassword string
		rotateGenerate bool
		jsonOutput     bool
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
	rotateCmd.Flags().BoolVarP(&rotateGenerate, "generate", "g", false, "Generate password")
	rotateCmd.Flags().BoolVarP(&jsonOutput, "json-output", "", false, "Whether output should be json or interactive (default interactive)")
	rotateCmd.MarkFlagRequired("username")
	rotateCmd.MarkFlagsMutuallyExclusive("password", "generate")
	rotateCmd.MarkFlagsOneRequired("password", "generate")
	root.AddCommand(rotateCmd)
}

func addBackupCommand(root *cobra.Command, cliCommand *CLICommand) {
	var (
		backupFile string
		backupUser string
		jsonOutput bool
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
	backupCmd.Flags().BoolVarP(&jsonOutput, "json-output", "", false, "Whether output should be json or interactive (default interactive)")
	backupCmd.MarkFlagRequired("file")
	root.AddCommand(backupCmd)
}
