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
}

type apiResponse struct {
	Formulae []struct {
		Name      string `json:"name"`
		Installed []struct {
			Time               int64 `json:"time"`
			InstalledOnRequest bool  `json:"installed_on_request"`
		} `json:"installed"`
	} `json:"formulae"`
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

	var pkgs []Package
	for _, f := range resp.Formulae {
		if len(f.Installed) == 0 {
			continue
		}
		install := f.Installed[0]
		if !install.InstalledOnRequest {
			continue
		}
		pkgs = append(pkgs, Package{
			Name:          f.Name,
			InstalledDate: time.Unix(install.Time, 0).Format("2006-01-02"),
		})
	}
	return pkgs, nil
}

func Uninstall(pkg string) error {
	return exec.Command("brew", "uninstall", pkg).Run()
}
