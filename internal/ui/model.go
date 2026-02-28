package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mrSamDev/brew-ui-potato/internal/brew"
)

const (
	rowNone         = ""
	rowUninstalling = "uninstalling"
	rowDeleted      = "deleted"

	colWidthPackage   = 32
	colWidthInstalled = 14
	colWidthStatus    = 14

	initialHeight = 20
	// heightOffset accounts for title, border, and footer rows.
	heightOffset = 7
)

type uninstallDoneMsg struct {
	idx int
	err error
}

// Model is the root Bubble Tea model.
type Model struct {
	table      table.Model
	packages   []brew.Package
	rowStatus  []string
	spinner    spinner.Model
	isLoading  bool
	loadingPkg string
	err        error
}

// InitialModel loads brew data and returns the starting model.
func InitialModel() Model {
	pkgs, err := brew.FetchPackages()

	rowStatus := make([]string, len(pkgs))

	columns := []table.Column{
		{Title: "Package", Width: colWidthPackage},
		{Title: "Installed", Width: colWidthInstalled},
		{Title: "Status", Width: colWidthStatus},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(buildRows(pkgs, rowStatus)),
		table.WithFocused(true),
		table.WithHeight(initialHeight),
		table.WithStyles(tableStyles),
	)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		table:     t,
		packages:  pkgs,
		rowStatus: rowStatus,
		spinner:   s,
		err:       err,
	}
}

// buildRows returns display rows with per-row status styles applied.
func buildRows(pkgs []brew.Package, status []string) []table.Row {
	rows := make([]table.Row, len(pkgs))
	for i, p := range pkgs {
		switch status[i] {
		case rowUninstalling:
			rows[i] = table.Row{
				redStyle.Render(p.Name),
				redStyle.Render(p.InstalledDate),
				redStyle.Render("Uninstalling..."),
			}
		case rowDeleted:
			rows[i] = table.Row{
				redStyle.Render(p.Name),
				redStyle.Render(p.InstalledDate),
				redStyle.Render("Deleted"),
			}
		default:
			rows[i] = table.Row{p.Name, p.InstalledDate, "User"}
		}
	}
	return rows
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		if k, ok := msg.(tea.KeyMsg); ok && (k.String() == "q" || k.String() == "ctrl+c") {
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - heightOffset)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case uninstallDoneMsg:
		return m.onUninstallDone(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) onUninstallDone(msg uninstallDoneMsg) (tea.Model, tea.Cmd) {
	m.isLoading = false
	m.loadingPkg = ""
	if msg.err != nil {
		m.rowStatus[msg.idx] = rowNone // revert on failure
	} else {
		m.rowStatus[msg.idx] = rowDeleted
	}
	m.table.SetRows(buildRows(m.packages, m.rowStatus))
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "d":
		if m.isLoading || m.rowStatus[m.table.Cursor()] != rowNone {
			return m, nil
		}
		return m.startUninstall(m.table.Cursor())
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
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
