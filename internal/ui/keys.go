package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/mrSamDev/brew-potato/internal/brew"
)

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.isConfirming {
		return m.handleConfirmKey(msg)
	}
	if m.isShowingAbout {
		return m.handleAboutKey(msg)
	}
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "?":
		m.isShowingAbout = true
		return m, nil
	case "d":
		if m.isLoading || m.isInitialLoading || len(m.packages) == 0 || m.rowStatus[m.table.Cursor()] != rowNone {
			return m, nil
		}
		m.isConfirming = true
		m.confirmIdx = m.table.Cursor()
		return m, nil
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) handleAboutKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc", "?", "b":
		m.isShowingAbout = false
		return m, nil
	}
	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		m.isConfirming = false
		return m.startUninstall(m.confirmIdx)
	case "n", "esc", "q", "ctrl+c":
		m.isConfirming = false
		return m, nil
	}
	return m, nil
}

func (m Model) startUninstall(idx int) (tea.Model, tea.Cmd) {
	pkg := m.packages[idx].Name
	m.isLoading = true
	m.loadingPkg = pkg
	m.rowStatus[idx] = rowUninstalling
	m.table.SetRows(buildRows(m.packages, m.rowStatus))
	return m, tea.Batch(
		func() tea.Msg {
			return uninstallDoneMsg{idx: idx, err: brew.Uninstall(pkg)}
		},
		m.spinner.Tick,
	)
}
