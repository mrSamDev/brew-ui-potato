package brew

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// Package represents a user-installed Homebrew formula.
type Package struct {
	Name          string
	InstalledDate string
	Type          string
}

const installedDateFormat = "2006-01-02"

type cask struct {
	Token         string `json:"token"`
	Installed     string `json:"installed"`
	InstalledTime int64  `json:"installed_time"`
}

type installEntry struct {
	Time               int64 `json:"time"`
	InstalledOnRequest bool  `json:"installed_on_request"`
}

type formula struct {
	Name      string         `json:"name"`
	Installed []installEntry `json:"installed"`
}

type apiResponse struct {
	Formulae []formula `json:"formulae"`
	Casks    []cask    `json:"casks"`
}

// FetchPackages returns all user-installed Homebrew formulae.
func FetchPackages() ([]Package, error) {
	out, err := exec.Command("brew", "info", "--json=v2", "--installed").Output()

	if err != nil {
		return nil, fmt.Errorf("brew info: %w", err)
	}

	var resp apiResponse
	if err = json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("parse brew output: %w", err)
	}

	pkgs := filterFormulaeRequest(resp.Formulae)
	caskPkgs := filterCasksRequest(resp.Casks)

	if len(pkgs) == 0 {
		// fall back to all installed formulae when installed_on_request is not
		// set (older homebrew or packages installed via scripts)
		pkgs = filterAllInstalled(resp.Formulae)
	}

	pkgs = append(pkgs, caskPkgs...)

	return pkgs, nil
}

func filterFormulaeRequest(formulae []formula) []Package {
	var pkgs []Package
	for _, f := range formulae {

		if len(f.Installed) == 0 {
			continue
		}

		if !f.Installed[0].InstalledOnRequest {
			continue
		}
		pkgs = append(pkgs, Package{
			Name:          f.Name,
			InstalledDate: time.Unix(f.Installed[0].Time, 0).Format(installedDateFormat),
			Type:          "formula",
		})
	}

	return pkgs
}

func filterCasksRequest(formulae []cask) []Package {
	var pkgs []Package
	for _, f := range formulae {
		if f.Installed == "" {
			continue
		}

		pkgs = append(pkgs, Package{
			Name:          f.Token,
			InstalledDate: time.Unix(f.InstalledTime, 0).Format(installedDateFormat),
			Type:          "cask",
		})
	}

	return pkgs
}

func filterAllInstalled(formulae []formula) []Package {
	var pkgs []Package
	for _, f := range formulae {
		if len(f.Installed) == 0 {
			continue
		}
		pkgs = append(pkgs, Package{
			Name:          f.Name,
			InstalledDate: time.Unix(f.Installed[0].Time, 0).Format(installedDateFormat),
			Type:          "formula",
		})
	}
	return pkgs
}

func Uninstall(pkg string) error {
	return exec.Command("brew", "uninstall", pkg).Run()
}
