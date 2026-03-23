package display

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Number of items to show per page in lists
// The Update function allows single digits for list row selection,
// so pageSize > 9 will cause issues with elements beyond row 9.
var pageSize = 9

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
	case *ConfirmState:
		return m.ConfirmView()
	case *EndState:
		return m.EndView()
	case *EmptyState:
		return m.EmptyView()
	default:
		return "Unknown state"
	}
}

func header() string {
	bannerStyle := lipgloss.NewStyle().
		Foreground((lipgloss.Color("117"))) // light blue
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("185")). // yellow
		Bold(true)
	msg := bannerStyle.Render(
		`
██╗  ██╗ █████╗  ██████╗ ██████╗
██║  ██║██╔══██╗██╔════╝██╔════╝
███████║███████║██║     ██║
██╔══██║██╔══██║██║     ██║
██║  ██║██║  ██║╚██████╗╚██████╗
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ `+
			versionStyle.Render("v2.0"),
	) + "\n"
	return msg + "\n" + horizontalLine() + "\n\n"
}

func (m model) WelcomeView() string {
	numServices := strconv.Itoa(m.exec.GetNumServices())
	msg := header() + welcomeFooterStyle.Render(
		"Vault contains "+numServices+" services.\nPress any key to continue...",
	)
	return msg
}

func (m model) EmptyView() string {
	return "The Vault is empty. Add a credential to get started!"
}

func (m model) EndView() string {
	var parts []string

	// 1. Render credential data first
	if len(m.result.Data) > 0 {
		for _, cred := range m.result.Data {
			line := fmt.Sprintf(
				"Service: %s | Username: %s",
				cred.Service,
				cred.Username,
			)

			if cred.Password != "" {
				line += fmt.Sprintf(" | Password: %s", cred.Password)
			}

			parts = append(parts, line)
		}
	}

	// 2. Render action result message
	if m.result.Success {
		parts = append(parts,
			fmt.Sprintf("%s executed %s command",
				successfully,
				m.result.Action,
			),
		)
	} else {
		parts = append(parts,
			fmt.Sprintf("%s executing %s command",
				failed,
				m.result.Action,
			),
		)
	}

	// 3. Append error message if present
	if m.result.Error != "" {
		parts = append(parts,
			fmt.Sprintf("Error: %s", m.result.Error),
		)
	}

	result := strings.Join(parts, "\n")
	return header() + addFooter(endStyle.Render(result), "")
}

func (m model) CredentialView() string {
	if len(m.result.Data) == 0 || m.result.Data[0].Service == "" {
		return "Error: no credential data available."
	}

	cred := m.result.Data[0]
	service := credServiceStyle.Render(" " + cred.Service + " ")
	user := credTextStyle.Render(cred.Username)
	pass := "*****"
	if m.showPass {
		pass = cred.Password
	}
	pass = credTextStyle.Render(pass)
	content := user + credSpacer + pass
	paddedContent := sidePaddingStyle.Render(content)
	return header() + "\n" + addFooter(credBox(service, paddedContent), credentialFooterStr)
}

func (m model) ConfirmView() string {
	service := credServiceStyle.Render(" " + m.cmd.Service + " ")
	user := credTextStyle.Render(m.cmd.Username)
	pass := credTextStyle.Render(m.cmd.Password)
	content := user + credSpacer + pass
	paddedContent := sidePaddingStyle.Render(content)
	return header() + credBox(service, paddedContent) + "\n" + footerStyle.Render("Use this password? y/n (default y)")
}

func (m model) ServiceListView() string {
	services := m.exec.GetServicesWithPrefix(m.cmd.Service)
	if len(services) == 0 {
		return header() + addFooter(fmt.Sprintf("No services matching %s.", m.cmd.Service), defaultFooterStr)
	}
	displayed_services := services[m.page*m.pageSize : min((m.page+1)*m.pageSize, len(services))]

	var rows [][]string
	for idx, serviceName := range displayed_services {
		rows = append(rows, []string{fmt.Sprintf("%d", idx+1), serviceName})
	}

	return header() + addFooter(
		listTable(
			"Service",
			rows,
			m.page+1,
			NumPages(len(services), m.pageSize),
			m.cursor,
		),
		defaultFooterStr,
	)
}

func (m model) UsernameListView() string {
	if m.cmd.Service == "" {
		return "Error: service name should not be empty in UsernameListView"
	}
	usernames, err := m.exec.GetUsersForServiceWithPrefix(m.cmd.Service, m.cmd.Username)
	if err != nil {
		return fmt.Sprintf("Error retrieving users for service %s: %v", m.cmd.Service, err)
	}
	if len(usernames) == 0 {
		return header() + addFooter(fmt.Sprintf("No usernames for service %s.", m.cmd.Service), defaultFooterStr)
	}

	displayed_usernames := usernames[m.page*m.pageSize : min((m.page+1)*m.pageSize, len(usernames))]
	var rows [][]string
	for idx, userName := range displayed_usernames {
		rows = append(rows, []string{fmt.Sprintf("%d", idx+1), userName})
	}
	headerRow := fmt.Sprintf("Usernames — %s", m.cmd.Service)
	if len(usernames) < 2 {
		headerRow = fmt.Sprintf("Username — %s", m.cmd.Service)
	}

	return header() + addFooter(
		listTable(
			headerRow,
			rows,
			m.page+1,
			NumPages(len(usernames), m.pageSize),
			m.cursor,
		),
		defaultFooterStr,
	)
}
