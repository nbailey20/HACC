package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
)

var (
	successTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42"))
	errorTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
	successfully = successTextStyle.Render("Successfully")
	failed       = errorTextStyle.Render("Failed")

	welcomeFooterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("141")). // light purple
				Bold(true)
	footerStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("244")). // grey
			Align(lipgloss.Center)

	defaultFooterStr = `
←→↑↓ / 1-9: navigate   Enter: select
Backspace: back        Esc / ^C: quit
`
	credentialFooterStr = `
Password copied.
Space: show/hide     Esc / ^C: quit
`

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("185")). // yellow
				Align(lipgloss.Center).
				Padding(0, 1)
	tableIndexStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("160")). // red
			Align(lipgloss.Center).
			Padding(0, 1)
	tableCursorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("117")). // light blue
				Bold(true).
				Foreground(lipgloss.Color("15")) // white
	tableRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")). // light blue
			Padding(0, 1)
	outerBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.Border{}). // invisible border to align pageText with table
				Padding(0, 1).
				Align(lipgloss.Center)
	pageNumberStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42")). // green
			Align(lipgloss.Center)

	credServiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("185"))               // yellow
	credTextStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("117"))               // light blue
	credSpacer       = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render(" │ ") // red
	sidePaddingStyle = lipgloss.NewStyle().Padding(0, 1)

	endStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("141")). // light purple
			Padding(0, 0).
			Align(lipgloss.Left)
)

func horizontalLine() string {
	width, _, _ := term.GetSize(os.Stdout.Fd())
	line := strings.Repeat("─", width)
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	return style.Render(line)
}

func addFooter(contentStr string, footerStr string) string {
	if footerStr == "" {
		footerStr = defaultFooterStr
	}
	footer := footerStyle.Render(footerStr)

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(contentStr + "\n" + footer)
}

func credBox(service string, content string) string {
	topLeft := "┌"
	topRight := "┐"
	bottomLeft := "└"
	bottomRight := "┘"
	horizontal := "─"
	vertical := "│"
	space := " "

	titleWidth := lipgloss.Width(service)
	contentWidth := max(lipgloss.Width(content), titleWidth)
	leftLineWidth := (contentWidth - titleWidth) / 2
	rightLineWidth := contentWidth - titleWidth - leftLineWidth
	contentPadding := max(0, (contentWidth-lipgloss.Width(content))/2)

	// build the box
	topContent := topLeft +
		strings.Repeat(horizontal, leftLineWidth) +
		service +
		strings.Repeat(horizontal, rightLineWidth) +
		topRight
	middleContent := vertical +
		strings.Repeat(space, contentPadding) +
		content +
		strings.Repeat(space, contentPadding) +
		vertical
	bottomContent := bottomLeft +
		strings.Repeat(horizontal, contentWidth) +
		bottomRight

	boxStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("141")) // light purple
	return boxStyle.Render(topContent + "\n" + middleContent + "\n" + bottomContent)
}

func listTable(header string, rows [][]string, pageNum int, totalPages int, cursor int) string {
	pageText := pageNumberStyle.
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
