package cli

type ActionKind int

const (
	ActionAdd ActionKind = iota
	ActionDelete
	ActionRotate
	ActionSearch
	ActionBackup
	ActionConfigure
	ActionUpgrade
	ActionInstall
	ActionEradicate
	ActionHelp
)

type CLIAction interface {
	Kind() ActionKind
}

type AddAction struct{}

func (AddAction) Kind() ActionKind { return ActionAdd }

type DeleteAction struct{}

func (DeleteAction) Kind() ActionKind { return ActionDelete }

type RotateAction struct{}

func (RotateAction) Kind() ActionKind { return ActionRotate }

type SearchAction struct{}

func (SearchAction) Kind() ActionKind { return ActionSearch }

type InstallAction struct{}
type EradicateAction struct{}
type UpgradeAction struct{}

type BackupAction struct{}

func (BackupAction) Kind() ActionKind { return ActionBackup }

type ConfigureAction struct{}

func (ConfigureAction) Kind() ActionKind { return ActionConfigure }

type HelpAction struct{}

func (HelpAction) Kind() ActionKind { return ActionHelp }

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
