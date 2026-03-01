package ui

import "testing"

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
