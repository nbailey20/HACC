package display

// interface for ErrorMsg types shown to user in EndState
type ErrorMsg interface {
	ErrorValue() error
	DisplayValue() string
}

type AddFailedMsg struct {
	Error   error
	Display string
}

func (m AddFailedMsg) ErrorValue() error {
	return m.Error
}

func (m AddFailedMsg) DisplayValue() string {
	return m.Display
}

type DeleteErrorMsg struct {
	Error   error
	Display string
}

func (m DeleteErrorMsg) ErrorValue() error {
	return m.Error
}

func (m DeleteErrorMsg) DisplayValue() string {
	return m.Display
}

type RotateErrorMsg struct {
	Error   error
	Display string
}

func (m RotateErrorMsg) ErrorValue() error {
	return m.Error
}

func (m RotateErrorMsg) DisplayValue() string {
	return m.Display
}

type BackupErrorMsg struct {
	Error   error
	Display string
}

func (m BackupErrorMsg) ErrorValue() error {
	return m.Error
}

func (m BackupErrorMsg) DisplayValue() string {
	return m.Display
}

// interface for SuccessMsg types shown to user in EndState
type SuccessMsg interface {
	DisplayValue() string
}

type AddSuccessMsg struct {
	Display string
}

func (m AddSuccessMsg) DisplayValue() string {
	return m.Display
}

type DeleteSuccessMsg struct {
	Display string
}

func (m DeleteSuccessMsg) DisplayValue() string {
	return m.Display
}

type RotateSuccessMsg struct {
	Display string
}

func (m RotateSuccessMsg) DisplayValue() string {
	return m.Display
}

type BackupSuccessMsg struct {
	Display string
}

func (m BackupSuccessMsg) DisplayValue() string {
	return m.Display
}

// other message types
type PasswordGeneratedMsg struct {
	Password string
}

type PasswordLoadedMsg struct {
	Password string
}

type PasswordCopiedMsg struct{}
