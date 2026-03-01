package ui

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

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

	// dialogPaddingH is the horizontal padding inside the confirmation dialog border.
	dialogPaddingH = 2
)

type uninstallDoneMsg struct {
	idx int
	err error
}

type packagesLoadedMsg struct {
	packages []brew.Package
	err      error
}

// Model is the root Bubble Tea model.
type Model struct {
	table            table.Model
	packages         []brew.Package
	rowStatus        []string
	spinner          spinner.Model
	isInitialLoading bool
	isLoading        bool
	isConfirming     bool
	isShowingAbout   bool
	confirmIdx       int
	loadingPkg       string
	err              error
	width            int
}

// InitialModel returns the starting model; package fetching happens in Init.
func InitialModel() Model {
	columns := []table.Column{
		{Title: "Package", Width: colWidthPackage},
		{Title: "Installed", Width: colWidthInstalled},
		{Title: "Status", Width: colWidthStatus},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(initialHeight),
		table.WithStyles(tableStyles),
	)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		table:            t,
		spinner:          s,
		isInitialLoading: true,
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
	return tea.Batch(
		func() tea.Msg {
			pkgs, err := brew.FetchPackages()
			return packagesLoadedMsg{packages: pkgs, err: err}
		},
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		if k, ok := msg.(tea.KeyPressMsg); ok && (k.String() == "q" || k.String() == "ctrl+c") {
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.table.SetWidth(colWidthPackage + colWidthInstalled + colWidthStatus)
		m.table.SetHeight(msg.Height - heightOffset)
		return m, nil
	case packagesLoadedMsg:
		return m.onPackagesLoaded(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case uninstallDoneMsg:
		return m.onUninstallDone(msg)
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) onPackagesLoaded(msg packagesLoadedMsg) (tea.Model, tea.Cmd) {
	m.isInitialLoading = false
	if msg.err != nil {
		m.err = msg.err
		return m, nil
	}
	m.packages = msg.packages
	m.rowStatus = make([]string, len(msg.packages))
	m.table.SetRows(buildRows(m.packages, m.rowStatus))
	return m, nil
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
