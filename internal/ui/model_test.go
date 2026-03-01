package ui

import (
	"errors"
	"testing"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

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


func TestInit_returnsCmd(t *testing.T) {
	m := newTestModel(testPkgs)
	m.isInitialLoading = true
	if cmd := m.Init(); cmd == nil {
		t.Error("Init should return a non-nil fetch command")
	}
}

func TestUpdate_packagesLoaded(t *testing.T) {
	m := newTestModel(nil)
	m.isInitialLoading = true

	newM, cmd := m.Update(packagesLoadedMsg{packages: testPkgs})

	if cmd != nil {
		t.Error("packagesLoadedMsg should return nil cmd")
	}
	updated := mustModel(t, newM)
	if updated.isInitialLoading {
		t.Error("isInitialLoading should be false after packages loaded")
	}
	if len(updated.packages) != len(testPkgs) {
		t.Errorf("packages len = %d, want %d", len(updated.packages), len(testPkgs))
	}
}

func TestUpdate_packagesLoaded_error(t *testing.T) {
	m := newTestModel(nil)
	m.isInitialLoading = true

	newM, _ := m.Update(packagesLoadedMsg{err: errors.New("brew failed")})

	updated := mustModel(t, newM)
	if updated.isInitialLoading {
		t.Error("isInitialLoading should be false even on error")
	}
	if updated.err == nil {
		t.Error("err should be set when packagesLoadedMsg contains an error")
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
		msg  tea.KeyPressMsg
	}{
		{"q", tea.KeyPressMsg{Code: 'q'}},
		{"ctrl+c", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}},
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

	_, cmd := m.Update(tea.KeyPressMsg{Code: 'd'})
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

