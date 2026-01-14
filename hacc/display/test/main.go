package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
)

func HorizontalLine() string {
	width, _, _ := term.GetSize(os.Stdout.Fd())
	line := strings.Repeat("─", width)
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	return style.Render(line)
}

func printListView() {
	var display string
	bannerStyle := lipgloss.NewStyle().
		Foreground((lipgloss.Color("117"))) // light blue
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("185")). // yellow
		Bold(true)
	header := bannerStyle.Render(
		`
██╗  ██╗ █████╗  ██████╗ ██████╗
██║  ██║██╔══██╗██╔════╝██╔════╝
███████║███████║██║     ██║
██╔══██║██╔══██║██║     ██║
██║  ██║██║  ██║╚██████╗╚██████╗
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ` +
			versionStyle.Render("v2.0"),
	)
	tableHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("185")).Align(lipgloss.Center).Padding(0, 1) // yellow
	tableIndexStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("160")).Align(lipgloss.Center).Padding(0, 1)  // red
	tableRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Padding(0, 1)
	outerBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.Border{}).
		BorderForeground(lipgloss.Color("141")). // light purple
		Padding(0, 1).MarginLeft(4).Align(lipgloss.Center)
	pageNumber := lipgloss.NewStyle().
		Bold(true).                       // Underline(true).
		Foreground(lipgloss.Color("42")). // green
		Align(lipgloss.Center).Render("1 / 1\n")
	footer := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("244")).
		Align(lipgloss.Center).
		Render(
			"\nNavigate with ↑/↓ or 0-9, use Enter to select.\nPress ESC or ctrl-c to quit.\n",
		)

	table := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("141"))).
		Headers("#", "Username").
		Rows(
			[]string{"1", "test1@gmail.com"},
			[]string{"2", "test2@gmail.com"},
			[]string{"3", "test3@gmail.com"},
			[]string{"4", "test4@gmail.com"},
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == -1:
				return tableHeaderStyle
			case col == 0:
				return tableIndexStyle
			default:
				return tableRowStyle
			}
		}).Render()

	outerBorder := outerBorderStyle.Render(table + "\n" + pageNumber)

	display += header + "\n" + HorizontalLine() + "\n\n" + outerBorder + footer
	fmt.Println(display)
}

func printCredView() {
	bannerStyle := lipgloss.NewStyle().
		Foreground((lipgloss.Color("117"))) // light blue
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("185")). // yellow
		Bold(true)
	header := bannerStyle.Render(
		`
██╗  ██╗ █████╗  ██████╗ ██████╗
██║  ██║██╔══██╗██╔════╝██╔════╝
███████║███████║██║     ██║
██╔══██║██╔══██║██║     ██║
██║  ██║██║  ██║╚██████╗╚██████╗
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ` +
			versionStyle.Render("v2.0"),
	)

	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("117"))
	service := lipgloss.NewStyle().Foreground(lipgloss.Color("185")).Render(" gmail ")
	user := textStyle.Render("bob@bob.com")
	pass := textStyle.Render("WhereverPartlyAnnualAmazing")
	spacer := lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render(" │ ") //red
	content := user + spacer + pass
	footer := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("244")).
		Align(lipgloss.Center).
		Render(
			"\nPress ESC or ctrl-c to quit, backspace to go back.\n",
		)

	// Content inside box
	paddedContent := lipgloss.NewStyle().Padding(0, 1).Render(content)
	width := lipgloss.Width(paddedContent)

	// Compose top border with embedded title
	leftCorner := "┌"
	rightCorner := "┐"
	horizontal := "─"

	titleWidth := lipgloss.Width(service)
	// totalLineWidth := width - 2 // subtract corners

	leftLineWidth := (width - titleWidth) / 2
	rightLineWidth := width - titleWidth - leftLineWidth

	topContent := leftCorner +
		strings.Repeat(horizontal, leftLineWidth) +
		service +
		strings.Repeat(horizontal, rightLineWidth) +
		rightCorner

	middleContent := "│" + paddedContent + "│"

	// Bottom border
	bottomContent := "└" + strings.Repeat(horizontal, width) + "┘"

	boxStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("141")) // light purple
	box := boxStyle.Render(topContent + "\n" + middleContent + "\n" + bottomContent)

	display := header + "\n" + HorizontalLine() + "\n\n" + box + footer
	fmt.Println(display)
}

func main() {
	printListView()
	printCredView()
}
