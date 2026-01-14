package display

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Number of items to show per page in lists
// The Update function allows single digits for list row selection,
// so pageSize > 9 will cause issues with elements beyond row 9.
var pageSize = 5

func (m model) View() string {
	switch m.state.(type) {
	case *WelcomeState:
		return m.WelcomeView()
	case *ServiceListState:
		return m.ServiceListView()
	case *UsernameListState:
		return m.UsernameListView()
	case *CredentialState:
		return m.CredentialView()
	case *MessageState:
		return m.MessageView()
	case *EndState:
		return m.EndView()
	default:
		fmt.Println("Unknown state type")
		return "Unknown state"
	}
}

func (m model) WelcomeView() string {
	var (
		bannerStyle = lipgloss.NewStyle().
				Foreground((lipgloss.Color("117"))) // light blue 117
		versionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")). // green
				Bold(true)
		footerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("141")). // light purple
				Bold(true)
	)
	msg := bannerStyle.Render(
		`
██╗  ██╗ █████╗  ██████╗ ██████╗
██║  ██║██╔══██╗██╔════╝██╔════╝
███████║███████║██║     ██║
██╔══██║██╔══██║██║     ██║
██║  ██║██║  ██║╚██████╗╚██████╗
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ `+
			versionStyle.Render("v2.0"),
	) + "\n" +
		footerStyle.Render("Press any key to continue...")

	return msg
}

func (m model) CredentialView() string {
	if m.serviceName == "" {
		return "Error: service name should not be empty in CredentialView"
	}

	msg := ""
	inside_s := ""
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")). // purple-ish
		Padding(0, 0)

	inside_s += m.message + "\n"
	inside_s += fmt.Sprintf("Service: %s\n\n", m.serviceName)
	inside_s += fmt.Sprintf("Username: %s\n\n", m.username)
	if m.showSecret {
		value, err := m.vault.Get(m.serviceName, m.username)
		if err != nil {
			return fmt.Sprintf("Error retrieving secret: %v", err)
		}
		inside_s += fmt.Sprintf("Secret: %s\n", value)
	} else {
		inside_s += "Secret: [hidden]\n"
	}

	msg += style.Render(inside_s)

	footer := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244")).Render(
		"\nUse Backspace to view all credentials. Press ESC or ctlr-c to quit.\n",
	)
	msg += footer

	return msg
}

func (m model) MessageView() string {
	msg := ""
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")). // purple-ish
		Padding(0, 0)

	msg += style.Render(m.message)

	footer := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244")).Render(
		"\nUse Backspace to view all credentials. Press ESC or ctlr-c to quit.\n",
	)
	msg += footer

	return msg
}

func (m model) EndView() string {
	// we're going to quit after this view is displayed, so no footer instructions
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("141")). // light purple
		Padding(0, 0)
	successfulStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))
	// errorStyle := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("196"))

	// errorStyle.Render("Error")
	return style.Render(successfulStyle.Render("Successfully")+" saved credential for gmail") + "\n"
	// return style.Render(errorStyle.Render("Error")+" saving credential for gmail") + "\n"
	// return style.Render(m.message) + "\n"
}

func (m model) ServiceListView() string {
	services := m.vault.ListServices(m.serviceName)
	displayed_services := services[m.page*m.pageSize : min((m.page+1)*m.pageSize, len(services))]

	var b string

	// ---- TITLE ----
	// 	titleText := `
	// _  _ ____ _  _ _    ___
	// |  | |__| |  | |     |
	//  \/  |  | |__| |___  |
	// `
	// 	titleText := `
	// ██  ██  ▄▄▄  ▄▄ ▄▄ ▄▄   ▄▄▄▄▄▄
	// ██▄▄██ ██▀██ ██ ██ ██     ██
	//  ▀██▀  ██▀██ ▀███▀ ██▄▄▄  ██
	// `
	titleText := "Credential Vault"

	tableWidth := 5 + 25 + 3 // Row + Name + Secret + spacing/padding
	title := lipgloss.NewStyle().
		Bold(true).                       // Underline(true).
		Foreground(lipgloss.Color("42")). // green
		Align(lipgloss.Center).
		Width(tableWidth).
		Render(titleText)
	// b += title + "\n"

	tableHeader := fmt.Sprintf("%-5s %-25s", "Row", "Name")
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("185")) // yellow
	tableContent := headerStyle.Render(tableHeader) + "\n"

	for i, serviceName := range displayed_services {
		cursor := " "
		rowStyle := lipgloss.NewStyle()
		if i == m.cursor {
			cursor = ">"
			rowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("117")).Bold(true) // light blue
			// Foreground(lipgloss.Color("15"))   // white
			// Bold(true)
		} else if i%2 == 0 {
			rowStyle = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("117"))
		} else {
			rowStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("117"))
		}

		line := fmt.Sprintf("%s %-5d %-25s", cursor, i+1, serviceName)
		tableContent += rowStyle.Render(line) + "\n"
	}

	tableContent += "\n" + title

	// Wrap table with border and optionally center
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("141")). // light purple
		Padding(0, 1).MarginLeft(4)
	b += tableStyle.Render(tableContent)

	// ---- FOOTER ----
	footer := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("244")).
		Align(lipgloss.Center).
		Render(
			"\nNavigate with ↑/↓ or 0-9, use Enter to select.\nPress ESC or ctlr-c to quit.\n",
		)
	b += footer
	return b
}

func (m model) UsernameListView() string {
	if m.serviceName == "" {
		return "Error: service name should not be empty in UsernameListView"
	}
	service, ok := m.vault.Services[m.serviceName]
	if !ok {
		return fmt.Sprintf("Error retrieving service %s: not found", m.serviceName)
	}
	usernames := service.GetUsers("")
	displayed_usernames := usernames[m.page*m.pageSize : min((m.page+1)*m.pageSize, len(usernames))]

	var b string

	// ---- TITLE ----
	tableWidth := 5 + 25 + 15 + 5 // Row + Name + Secret + spacing/padding
	title := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("39")). // cyan
		Width(tableWidth).                // match table width
		Align(lipgloss.Center).
		Render(fmt.Sprintf("USERNAMES FOR %s", m.serviceName))
	b += "\n" + title + "\n"

	tableHeader := fmt.Sprintf("%-5s %-25s", "Row", "Username")
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")) // yellow
	tableContent := headerStyle.Render(tableHeader) + "\n"

	for i, username := range displayed_usernames {
		cursor := " "
		rowStyle := lipgloss.NewStyle()
		if i == m.cursor {
			cursor = ">"
			rowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("57")). // bright blue
				Foreground(lipgloss.Color("15")). // white
				Bold(true)
		} else if i%2 == 0 {
			rowStyle = lipgloss.NewStyle().Background(lipgloss.Color("235"))
		} else {
			rowStyle = lipgloss.NewStyle().Background(lipgloss.Color("236"))
		}

		line := fmt.Sprintf("%s %-5d %-25s", cursor, i+1, username)
		tableContent += rowStyle.Render(line) + "\n"
	}

	// Wrap table with border and optionally center
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)
	b += tableStyle.Render(tableContent)

	// ---- FOOTER ----
	footer := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244")).Render(
		"\nUse ↑/↓ or type row number, press Enter to view. Press ESC or ctlr-c to quit.\n",
	)
	b += footer
	return b
}
