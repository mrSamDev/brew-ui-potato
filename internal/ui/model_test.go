package ui

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mrSamDev/brew-ui-potato/internal/brew"
)

var testPkgs = []brew.Package{
	{Name: "git", InstalledDate: "2023-11-14"},
	{Name: "wget", InstalledDate: "2024-03-09"},
}

func newTestModel(pkgs []brew.Package) Model {
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
	}
}

func mustModel(t *testing.T, m tea.Model) Model {
	t.Helper()
	model, ok := m.(Model)
	if !ok {
		t.Fatalf("expected Model, got %T", m)
	}
	return model
}

func TestBuildRows_empty(t *testing.T) {
	rows := buildRows(nil, nil)
	if len(rows) != 0 {
		t.Errorf("got %d rows, want 0", len(rows))
	}
}

func TestBuildRows_defaultStatus(t *testing.T) {
	rows := buildRows(testPkgs[:1], []string{rowNone})

	if rows[0][0] != "git" {
		t.Errorf("name = %q, want %q", rows[0][0], "git")
	}
	if rows[0][2] != "User" {
		t.Errorf("status = %q, want %q", rows[0][2], "User")
	}
}

func TestBuildRows_uninstallingStatus(t *testing.T) {
	rows := buildRows(testPkgs[:1], []string{rowUninstalling})

	cell := rows[0][2]
	// lipgloss may strip ANSI in test environments; check plain text
	if cell != "Uninstalling..." && cell != redStyle.Render("Uninstalling...") {
		t.Errorf("status cell %q does not match expected uninstalling text", cell)
	}
}

func TestBuildRows_deletedStatus(t *testing.T) {
	rows := buildRows(testPkgs[:1], []string{rowDeleted})

	cell := rows[0][2]
	if cell != "Deleted" && cell != redStyle.Render("Deleted") {
		t.Errorf("status cell %q does not match expected deleted text", cell)
	}
}

func TestInit_returnsNil(t *testing.T) {
	m := newTestModel(testPkgs)
	if cmd := m.Init(); cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestUpdate_windowResize(t *testing.T) {
	m := newTestModel(testPkgs)
	newM, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

	if cmd != nil {
		t.Error("WindowSizeMsg should return nil cmd")
	}
	if _, ok := newM.(Model); !ok {
		t.Errorf("expected Model, got %T", newM)
	}
}

func TestUpdate_errorState_quitsOnQuitKeys(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyMsg
	}{
		{"q", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel(nil)
			m.err = errors.New("brew failed")

			_, cmd := m.Update(tt.msg)
			if cmd == nil {
				t.Fatal("expected non-nil cmd")
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Error("expected QuitMsg")
			}
		})
	}
}

func TestUpdate_errorState_ignoresOtherKeys(t *testing.T) {
	m := newTestModel(nil)
	m.err = errors.New("brew failed")

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if cmd != nil {
		t.Error("non-quit key in error state should return nil cmd")
	}
}

func TestUpdate_uninstallDone_success(t *testing.T) {
	m := newTestModel(testPkgs)
	m.rowStatus[0] = rowUninstalling
	m.isLoading = true
	m.loadingPkg = "git"

	newM, cmd := m.Update(uninstallDoneMsg{idx: 0, err: nil})

	if cmd != nil {
		t.Error("uninstallDoneMsg should return nil cmd")
	}
	updated := mustModel(t, newM)
	if updated.rowStatus[0] != rowDeleted {
		t.Errorf("rowStatus[0] = %q, want %q", updated.rowStatus[0], rowDeleted)
	}
	if updated.isLoading {
		t.Error("isLoading should be false after uninstall completes")
	}
	if updated.loadingPkg != "" {
		t.Errorf("loadingPkg = %q, want empty", updated.loadingPkg)
	}
}

func TestUpdate_uninstallDone_failure(t *testing.T) {
	m := newTestModel(testPkgs)
	m.rowStatus[0] = rowUninstalling
	m.isLoading = true

	newM, _ := m.Update(uninstallDoneMsg{idx: 0, err: errors.New("uninstall failed")})

	updated := mustModel(t, newM)
	if updated.rowStatus[0] != rowNone {
		t.Errorf("rowStatus[0] = %q, want %q on failure", updated.rowStatus[0], rowNone)
	}
	if updated.isLoading {
		t.Error("isLoading should be false after uninstall failure")
	}
}

func TestHandleKey_quitKeys(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyMsg
	}{
		{"q", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel(testPkgs)
			_, cmd := m.Update(tt.msg)

			if cmd == nil {
				t.Fatal("expected non-nil cmd")
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Error("expected QuitMsg")
			}
		})
	}
}

func TestHandleKey_delete_setsUninstalling(t *testing.T) {
	m := newTestModel(testPkgs)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	model := mustModel(t, updated)
	if !model.isLoading {
		t.Error("isLoading should be true after d key")
	}
	if model.rowStatus[0] != rowUninstalling {
		t.Errorf("rowStatus[0] = %q, want %q", model.rowStatus[0], rowUninstalling)
	}
	if cmd == nil {
		t.Error("d should return a non-nil command")
	}
}

func TestHandleKey_delete_skipsWhenLoading(t *testing.T) {
	m := newTestModel(testPkgs)
	m.isLoading = true

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if cmd != nil {
		t.Error("d while loading should return nil cmd")
	}
}

func TestHandleKey_delete_skipsWhenRowNotIdle(t *testing.T) {
	m := newTestModel(testPkgs)
	m.rowStatus[0] = rowUninstalling

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if cmd != nil {
		t.Error("d on non-idle row should return nil cmd")
	}
}
