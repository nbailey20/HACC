package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var pageSize = 5

type model struct {
	vault       Vault
	serviceName string
	page        int
	pageSize    int
	cursor      int
	message     string
	showSecret  bool
	state       State
}

func NewModel(cli CLIInput, vault Vault) *model {
	switch cli.Action {
	case "get":
		if cli.Name != "" {
			return &model{
				vault:       vault,
				serviceName: cli.Name,
				page:        0,
				pageSize:    pageSize,
				cursor:      0,
				message:     "",
				showSecret:  true,
				state:       &DetailState{},
			}
		}
		return &model{
			vault:       vault,
			serviceName: "",
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  true,
			state:       &WelcomeState{},
		}
	case "add":
		vault.Add(cli.Name, cli.Value)
		return &model{
			vault:       vault,
			serviceName: cli.Name,
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "Added the following service to the vault:\n",
			showSecret:  false,
			state:       &DetailState{},
		}
	case "delete":
		err := vault.Delete(cli.Name)
		msg := fmt.Sprintf("Service %s deleted successfully.", cli.Name)
		if err != nil {
			msg = fmt.Sprintf("Error deleting service %s: %v\n", cli.Name, err)
		}
		return &model{
			vault:       vault,
			serviceName: "",
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     msg,
			showSecret:  false,
			state:       &MessageState{},
		}
	default:
		return &model{
			vault:       vault,
			serviceName: "",
			page:        0,
			pageSize:    pageSize,
			cursor:      0,
			message:     "",
			showSecret:  false,
			state:       &WelcomeState{},
		}
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp:
			return m.state.Update(m, UpEvent{})
		case tea.KeyDown:
			return m.state.Update(m, DownEvent{})
		case tea.KeyLeft:
			return m.state.Update(m, LeftEvent{})
		case tea.KeyRight:
			return m.state.Update(m, RightEvent{})
		case tea.KeyEnter:
			return m.state.Update(m, EnterEvent{})
		case tea.KeyBackspace:
			return m.state.Update(m, BackEvent{})
		case tea.KeyRunes:
			// Numbers (and other printable characters)
			if len(msg.Runes) == 1 && msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
				n := int(msg.Runes[0] - '0')
				return m.state.Update(m, NumberEvent{Number: n})
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state.(type) {
	case *WelcomeState:
		return m.WelcomeView()
	case *ListState:
		return m.ListView()
	case *DetailState:
		return m.DetailView()
	case *MessageState:
		return m.MessageView()
	default:
		fmt.Println("Unknown state type")
		return "Unknown state"
	}
}

func (m model) WelcomeView() string {
	return "Welcome to the Credential Vault!\nPress any key to continue..."
}

func (m model) ListView() string {
	if m.serviceName != "" {
		return fmt.Sprintf("Error: service name should be empty in ListView, got %s", m.serviceName)
	}
	services := m.vault.ListServices()
	displayed_services := services[m.page*m.pageSize : min((m.page+1)*m.pageSize, len(services))]

	var b string

	// ---- TITLE ----
	tableWidth := 5 + 25 + 15 + 5 // Row + Name + Secret + spacing/padding
	title := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("39")). // cyan
		Width(tableWidth).                // match table width
		Align(lipgloss.Center).
		Render("CREDENTIAL VAULT")
	b += "\n" + title + "\n"

	tableHeader := fmt.Sprintf("%-5s %-25s", "Row", "Name")
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")) // yellow
	tableContent := headerStyle.Render(tableHeader) + "\n"

	for i, serviceName := range displayed_services {
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

		line := fmt.Sprintf("%s %-5d %-25s", cursor, i+1, serviceName)
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

func (m model) DetailView() string {
	if m.serviceName == "" {
		return "Error: service name should not be empty in DetailView"
	}

	s := m.message + "\n"
	s += fmt.Sprintf("Service: %s\n\n", m.serviceName)
	if m.showSecret {
		value, err := m.vault.Get(m.serviceName)
		if err != nil {
			return fmt.Sprintf("Error retrieving secret: %v", err)
		}
		s += fmt.Sprintf("Secret: %s\n", value)
	} else {
		s += "Secret: [hidden]\n"
	}
	return s
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
