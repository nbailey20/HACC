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
	msg := header() + welcomeFooterStyle.Render("Press any key to continue...")
	return msg
}

func (m model) EmptyView() string {
	return "The Vault is empty. Add a credential to get started!"
}

func (m model) EndView() string {
	result := m.endMessage
	if !m.endSuccess {
		result += fmt.Sprintf("\nError: %v", m.endError)
	}
	// we're going to quit after this view is displayed, so no footer instructions
	return header() + addFooter(endStyle.Render(result), "")
}

func (m model) CredentialView() string {
	if m.serviceName == "" {
		return "Error: service name should not be empty in CredentialView"
	}

	service := credServiceStyle.Render(" " + m.serviceName + " ")
	user := credTextStyle.Render(m.username)
	pass := "*****"
	if m.password == "" {
		pass = "Loading..."
	}
	if m.showPass {
		pass = m.password
	}
	pass = credTextStyle.Render(pass)
	content := user + credSpacer + pass
	paddedContent := sidePaddingStyle.Render(content)
	return header() + "\n" + addFooter(credBox(service, paddedContent), credentialFooterStr)
}

func (m model) ConfirmView() string {
	service := credServiceStyle.Render(" " + m.serviceName + " ")
	user := credTextStyle.Render(m.username)
	pass := credTextStyle.Render(m.password)
	content := user + credSpacer + pass
	paddedContent := sidePaddingStyle.Render(content)
	return header() + credBox(service, paddedContent) + "\n" + footerStyle.Render("Use this password? y/n (default y)")
}

func (m model) ServiceListView() string {
	services := m.vault.ListServices(m.serviceName)
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
	if m.serviceName == "" {
		return "Error: service name should not be empty in UsernameListView"
	}

	usernames, err := m.vault.GetUsersForService(m.serviceName)
	if err != nil {
		return fmt.Sprintf("Error retrieving users for service %s: %v", m.serviceName, err)
	}
	displayed_usernames := usernames[m.page*m.pageSize : min((m.page+1)*m.pageSize, len(usernames))]
	var rows [][]string
	for idx, userName := range displayed_usernames {
		rows = append(rows, []string{fmt.Sprintf("%d", idx+1), userName})
	}

	return header() + addFooter(
		listTable(
			"Username",
			rows,
			m.page+1,
			NumPages(len(usernames), m.pageSize),
			m.cursor,
		),
		defaultFooterStr,
	)
}
