package display

import "github.com/charmbracelet/lipgloss"

var (
	successTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42"))

	errorTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	footerStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("244")). // grey
			Align(lipgloss.Center)

	successfully = successTextStyle.Render("Successfully")
	failed       = errorTextStyle.Render("Failed")

	defaultFooterStr = `
←→↑↓ / 1-9: navigate   Enter: select
Backspace: back        Esc / ^C: quit
`
	credentialFooterStr = `
Password copied.
Space: show/hide     Esc / ^C: quit
`
)
