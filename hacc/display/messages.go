package display

import "github.com/nbailey20/hacc/engine"

type resultMsg struct {
	result engine.Result
}

// other message types
type PasswordGeneratedMsg struct {
	Password string
}

type PasswordLoadedMsg struct {
	Password string
}

type PasswordCopiedMsg struct{}
