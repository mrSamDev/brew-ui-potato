package brew_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mrSamDev/brew-ui-potato/internal/brew"
)

// brewJSONMultiple contains four formulae that exercise all filter paths:
// git and wget are installed on request, dep-only is a dependency, curl has no installs.
const brewJSONMultiple = `{
	"formulae": [
		{
			"name": "git",
			"installed": [{"time": 1700000000, "installed_on_request": true}]
		},
		{
			"name": "dep-only",
			"installed": [{"time": 1700000000, "installed_on_request": false}]
		},
		{
			"name": "curl",
			"installed": []
		},
		{
			"name": "wget",
			"installed": [{"time": 1710000000, "installed_on_request": true}]
		}
	]
}`

// setupFakeBrew writes a shell script that acts as brew, returning the given JSON output.
func setupFakeBrew(t *testing.T, output string) {
	t.Helper()
	dir := t.TempDir()

	jsonFile := filepath.Join(dir, "output.json")
	if err := os.WriteFile(jsonFile, []byte(output), 0644); err != nil {
		t.Fatalf("write json file: %v", err)
	}

	script := "#!/bin/sh\ncat " + jsonFile + "\n"
	brewPath := filepath.Join(dir, "brew")
	if err := os.WriteFile(brewPath, []byte(script), 0755); err != nil {
		t.Fatalf("write brew script: %v", err)
	}

	t.Setenv("PATH", dir+":"+os.Getenv("PATH"))
}

// setupFailingBrew places a brew stub that always exits with code 1.
func setupFailingBrew(t *testing.T) {
	t.Helper()
	dir := t.TempDir()

	brewPath := filepath.Join(dir, "brew")
	if err := os.WriteFile(brewPath, []byte("#!/bin/sh\nexit 1\n"), 0755); err != nil {
		t.Fatalf("write brew script: %v", err)
	}

	t.Setenv("PATH", dir+":"+os.Getenv("PATH"))
}

func TestFetchPackages_returnsOnlyInstalledOnRequest(t *testing.T) {
	setupFakeBrew(t, brewJSONMultiple)

	pkgs, err := brew.FetchPackages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("got %d packages, want 2", len(pkgs))
	}
	if pkgs[0].Name != "git" {
		t.Errorf("pkgs[0].Name = %q, want %q", pkgs[0].Name, "git")
	}
	if pkgs[1].Name != "wget" {
		t.Errorf("pkgs[1].Name = %q, want %q", pkgs[1].Name, "wget")
	}
}

func TestFetchPackages_formatsDateAsYYYYMMDD(t *testing.T) {
	setupFakeBrew(t, brewJSONMultiple)

	pkgs, err := brew.FetchPackages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Unix(1700000000, 0).Format("2006-01-02")
	if pkgs[0].InstalledDate != want {
		t.Errorf("InstalledDate = %q, want %q", pkgs[0].InstalledDate, want)
	}
}

func TestFetchPackages_emptyFormulae(t *testing.T) {
	setupFakeBrew(t, `{"formulae":[]}`)

	pkgs, err := brew.FetchPackages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 0 {
		t.Errorf("got %d packages, want 0", len(pkgs))
	}
}

func TestFetchPackages_invalidJSON(t *testing.T) {
	setupFakeBrew(t, "not valid json")

	_, err := brew.FetchPackages()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestFetchPackages_brewCommandFails(t *testing.T) {
	setupFailingBrew(t)

	_, err := brew.FetchPackages()
	if err == nil {
		t.Fatal("expected error when brew exits non-zero, got nil")
	}
}

func TestUninstall_success(t *testing.T) {
	setupFakeBrew(t, `{"formulae":[]}`)

	if err := brew.Uninstall("git"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUninstall_failure(t *testing.T) {
	setupFailingBrew(t)

	if err := brew.Uninstall("git"); err == nil {
		t.Fatal("expected error when brew fails, got nil")
	}
}
