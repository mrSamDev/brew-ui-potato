package ui

import (
	"errors"
	"strings"
	"testing"
)

func TestView_withError(t *testing.T) {
	m := newTestModel(nil)
	m.err = errors.New("command not found: brew")

	view := m.View()

	if !strings.Contains(view, "Failed to load brew data:") {
		t.Error("error view should contain error header")
	}
	if !strings.Contains(view, "command not found: brew") {
		t.Error("error view should contain the error message")
	}
	if !strings.Contains(view, "Press q to quit") {
		t.Error("error view should contain quit instruction")
	}
}

func TestView_normalState(t *testing.T) {
	m := newTestModel(testPkgs)

	view := m.View()

	if !strings.Contains(view, "Brew Packages") {
		t.Error("normal view should contain the title")
	}
	if !strings.Contains(view, "↑/↓  navigate") {
		t.Error("normal view should contain navigation hint")
	}
	if !strings.Contains(view, "d  uninstall") {
		t.Error("normal view should contain uninstall hint")
	}
}

func TestView_loadingState(t *testing.T) {
	m := newTestModel(testPkgs)
	m.isLoading = true
	m.loadingPkg = "testpkg"

	view := m.View()

	if !strings.Contains(view, "Uninstalling testpkg") {
		t.Errorf("loading view should contain package name, got:\n%s", view)
	}
	// Navigation hints should not appear while loading
	if strings.Contains(view, "↑/↓  navigate") {
		t.Error("loading view should not show navigation hint")
	}
}
