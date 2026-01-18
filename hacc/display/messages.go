package display

type AddFailedMsg struct {
	Error error
}

type DeleteFailedMsg struct {
	Error error
}

type RotateFailedMsg struct {
	Error error
}

type PasswordGeneratedMsg struct {
	Password string
}

type PasswordGenerationErrorMsg struct {
	Error error
}

type BackupErrorMsg struct {
	Error error
}
