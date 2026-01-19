package display

import "github.com/charmbracelet/lipgloss"

var (
	successTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42"))

	errorTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	successfully = successTextStyle.Render("Successfully")
	failed       = errorTextStyle.Render("Failed")
)
