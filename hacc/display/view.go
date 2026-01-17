package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
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

func horizontalLine() string {
	width, _, _ := term.GetSize(os.Stdout.Fd())
	line := strings.Repeat("─", width)
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	return style.Render(line)
}

func addFooter(content string) string {
	rawFooter := `
←→↑↓ / 1-9: navigate   Enter: select
Backspace: back        Esc / ^C: quit
`
	footer := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("244")). // grey
		Align(lipgloss.Center).
		Render(
			rawFooter,
		)

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(content + "\n" + footer)
}

func credBox(service string, content string) string {
	topLeft := "┌"
	topRight := "┐"
	bottomLeft := "└"
	bottomRight := "┘"
	horizontal := "─"
	vertical := "│"

	titleWidth := lipgloss.Width(service)
	contentWidth := lipgloss.Width(content)
	leftLineWidth := (contentWidth - titleWidth) / 2
	rightLineWidth := contentWidth - titleWidth - leftLineWidth

	// build the box
	topContent := topLeft +
		strings.Repeat(horizontal, leftLineWidth) +
		service +
		strings.Repeat(horizontal, rightLineWidth) +
		topRight
	middleContent := vertical + content + vertical
	bottomContent := bottomLeft + strings.Repeat(horizontal, contentWidth) + bottomRight

	boxStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("141")) // light purple
	return boxStyle.Render(topContent + "\n" + middleContent + "\n" + bottomContent)
}

func listTable(header string, rows [][]string, pageNum int, totalPages int, cursor int) string {
	tableHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("185")). // yellow
		Align(lipgloss.Center).
		Padding(0, 1)
	tableIndexStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("160")). // red
		Align(lipgloss.Center).
		Padding(0, 1)
	tableCursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("117")). // light blue
		Bold(true).
		Foreground(lipgloss.Color("15")) // white
	tableRowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")). // light blue
		Padding(0, 1)
	outerBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.Border{}). // invisible border to align pageText with table // BorderForeground(lipgloss.Color("141")). // light purple
		Padding(0, 1).
		Align(lipgloss.Center)
	pageText := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42")). // green
		Align(lipgloss.Center).
		Render(fmt.Sprintf("%d / %d", pageNum, totalPages))

	table := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("141"))). // purple
		Headers("#", header).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == -1:
				return tableHeaderStyle
			case col == 0:
				return tableIndexStyle
			case row == cursor:
				return tableCursorStyle
			default:
				return tableRowStyle
			}
		}).Render()

	return outerBorderStyle.Render(table + "\n" + pageText)
}

func (m model) WelcomeView() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")). // light purple
		Bold(true)
	msg := header() + footerStyle.Render("Press any key to continue...")
	return msg
}

func (m model) EmptyView() string {
	return "The Vault is empty. Add a credential to get started!"
}

func (m model) EndView() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("141")). // light purple
		Padding(0, 0)
	successText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Render("Success ")
	errorText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Render("Error ")

	result := successText
	if !m.endSuccess {
		result = errorText
	}
	switch m.action {
	case "add":
		result += "adding username " + m.username + " for " + m.serviceName + ".\n"
	case "delete":
		result += "deleting username " + m.username + " for " + m.serviceName + ".\n"
	case "rotate":
		result += "rotating username " + m.username + " for " + m.serviceName + ".\n"
	}
	// we're going to quit after this view is displayed, so no footer instructions
	return header() + style.Render(result) + "\n"
}

func (m model) CredentialView() string {
	if m.serviceName == "" {
		return "Error: service name should not be empty in CredentialView"
	}

	service := lipgloss.NewStyle().Foreground(lipgloss.Color("185")).Render(" gmail ") // yellow
	credTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("117"))             // light blue
	user := credTextStyle.Render(m.username)
	rawPass, err := m.vault.Get(m.serviceName, m.username)
	pass := credTextStyle.Render(rawPass)
	if err != nil {
		pass = credTextStyle.Render(fmt.Sprintf("%v", err))
		// return fmt.Sprintf("Error retrieving secret: %v", err)
	}
	spacer := lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render(" │ ") //red
	content := user + spacer + pass
	paddedContent := lipgloss.NewStyle().Padding(0, 1).Render(content)

	return header() + "\n" + addFooter(credBox(service, paddedContent))
}

func (m model) ServiceListView() string {
	services := m.vault.ListServices(m.serviceName)
	displayed_services := services[m.page*m.pageSize : min((m.page+1)*m.pageSize, len(services))]

	var rows [][]string
	for idx, serviceName := range displayed_services {
		rows = append(rows, []string{fmt.Sprintf("%d", idx+1), serviceName})
	}

	return header() + addFooter(listTable(
		"Service",
		rows,
		m.page+1,
		NumPages(len(services), m.pageSize),
		m.cursor,
	))
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

	var rows [][]string
	for idx, userName := range displayed_usernames {
		rows = append(rows, []string{fmt.Sprintf("%d", idx+1), userName})
	}

	return header() + addFooter(listTable(
		"Username",
		rows,
		m.page+1,
		NumPages(len(usernames), m.pageSize),
		m.cursor,
	))
}
