package cli

// func NormalizeInput(input CLIInput) RawCLICommand {
// 	var action []string
// 	if input.Install {
// 		action = append(action, "install")
// 	}
// 	if input.Eradicate {
// 		action = append(action, "eradicate")
// 	}
// 	if input.Add {
// 		action = append(action, "add")
// 	}
// 	if input.Delete {
// 		action = append(action, "delete")
// 	}
// 	if input.Rotate {
// 		action = append(action, "rotate")
// 	}
// 	if input.Backup {
// 		action = append(action, "backup")
// 	}
// 	if input.Configure {
// 		action = append(action, "configure")
// 	}
// 	if input.Upgrade {
// 		action = append(action, "upgrade")
// 	}

// 	if len(action) == 0 {
// 		action = append(action, "search") // default action
// 	}

// 	return RawCLICommand{
// 		Action:   action,
// 		Service:  input.Service,
// 		Username: input.Username,
// 		Password: input.Password,
// 		Generate: input.Generate,
// 		File:     input.File,
// 		Wipe:     input.Wipe,
// 		Export:   input.Export,
// 		Set:      input.Set,
// 		Show:     input.Show,
// 	}
// }
