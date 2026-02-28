package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n  %s\n\n  %s\n\n  Press q to quit.\n",
			errorStyle.Render("Failed to load brew data:"),
			m.err.Error(),
		)
	}

	title := titleStyle.Render("🍺  Brew Packages")
	tableView := tableWrapStyle.Render(m.table.View())

	var footer string
	if m.isLoading {
		footer = footerStyle.Render(
			fmt.Sprintf("  %s  Uninstalling %s...", m.spinner.View(), m.loadingPkg),
		)
	} else {
		footer = footerStyle.Render("  ↑/↓  navigate    d  uninstall    q  quit")
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, tableView, footer)
}
