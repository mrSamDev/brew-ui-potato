package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorAccent = lipgloss.Color("99")  // bright purple
	colorBorder = lipgloss.Color("62")  // muted purple
	colorMuted  = lipgloss.Color("241") // dim grey

	redStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottom(true).
			BorderForeground(colorBorder).
			Padding(0, 1).
			MarginBottom(1)

	tableWrapStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	footerStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	tableStyles = func() table.Styles {
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			BorderBottom(true).
			Bold(true).
			Foreground(colorAccent)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("230")).
			Background(colorBorder).
			Bold(true)
		return s
	}()
)
