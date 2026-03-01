package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestHandleKey_quitKeys(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
	}{
		{"q", tea.KeyPressMsg{Code: 'q'}},
		{"ctrl+c", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}},
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

func TestHandleKey_delete_opensConfirmation(t *testing.T) {
	m := newTestModel(testPkgs)

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'd'})

	model := mustModel(t, updated)
	if !model.isConfirming {
		t.Error("isConfirming should be true after d key")
	}
	if model.confirmIdx != 0 {
		t.Errorf("confirmIdx = %d, want 0", model.confirmIdx)
	}
	if cmd != nil {
		t.Error("d should return a nil cmd — uninstall starts only after confirmation")
	}
}

func TestHandleKey_delete_skipsWhenLoading(t *testing.T) {
	m := newTestModel(testPkgs)
	m.isLoading = true

	_, cmd := m.Update(tea.KeyPressMsg{Code: 'd'})
	if cmd != nil {
		t.Error("d while loading should return nil cmd")
	}
}

func TestHandleKey_delete_skipsWhenRowNotIdle(t *testing.T) {
	m := newTestModel(testPkgs)
	m.rowStatus[0] = rowUninstalling

	_, cmd := m.Update(tea.KeyPressMsg{Code: 'd'})
	if cmd != nil {
		t.Error("d on non-idle row should return nil cmd")
	}
}
