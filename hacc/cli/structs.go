package cli

type CLIAction interface{}

type AddAction struct{}
type DeleteAction struct{}
type RotateAction struct{}
type SearchAction struct{}
type InstallAction struct{}
type EradicateAction struct{}
type BackupAction struct{}
type ConfigureAction struct{}
type UpgradeAction struct{}

type CLICommand struct {
	Action   CLIAction
	Service  string
	Username string
	Password string
	Generate bool
	File     string
	Wipe     bool
	Export   bool
	Set      string
	Show     string
}
