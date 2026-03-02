# brew-potato

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/license-do%20what%20you%20want-875fff)](LICENSE)
[![GitHub](https://img.shields.io/badge/github-brew--potato-181717?logo=github)](https://github.com/mrSamDev/brew-potato)

A terminal UI for managing Homebrew packages, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Why

`brew list` dumps everything: your packages, their dependencies, transitive dependencies, all tangled together. There's no built-in way to see what _you_ actually installed.

I knew the alternatives. Nix exists. There's even a Homebrew rewrite in Rust (zerobrew). I tried them, forgot about them, and went back to `brew`.

So I built this: show only user-installed packages, nothing else. I was also learning Go and needed a real project, not another todo app.

## What it does

Full-screen TUI with a rounded table border and styled header and footer. Lists your user-installed Homebrew formulae with install dates. Uninstall runs async with a live spinner; the row turns red in place, no view flicker.

If `brew` isn't available, it tells you cleanly instead of dumping a stack trace on you.

## Built With

| Library | Purpose |
| ------- | ------- |
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [Bubbles](https://github.com/charmbracelet/bubbles) | Table and spinner components |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | Styles and layout |

![brew-potato UI](https://res.cloudinary.com/dnmuyrcd7/image/upload/UI_ofnace.png)

## Requirements

- [Go](https://go.dev/) 1.21+
- [Homebrew](https://brew.sh/)

## Run

```sh
go mod tidy
go run .
```

## Build

```sh
go build -o brew-potato .
./brew-potato
```

## Project Structure

```
brew-ui/
├── main.go                 # entry point
├── internal/
│   ├── brew/
│   │   └── brew.go         # brew CLI wrapper & data types
│   └── ui/
│       ├── model.go        # Bubble Tea model, Init, Update
│       ├── view.go         # View rendering
│       └── styles.go       # lipgloss styles & table theme
├── go.mod
└── go.sum
```

## Release With Homebrew (GoReleaser)

This repo is configured for GoReleaser v2 (`.goreleaser.yaml`) and GitHub Actions (`.github/workflows/release.yml`) to:

- build release binaries for macOS and Linux
- create GitHub release assets + checksums
- publish a Homebrew cask into a tap repo (`mrSamDev/homebrew-tap`)

### One-time setup

1. Create the tap repo (or change `.goreleaser.yaml` to your own):
   - `mrSamDev/homebrew-tap`


### Create a release

Version comes from the Git tag you push (for example `v0.1.0`):

```sh
git tag v0.1.0
git push origin v0.1.0
```

That tag triggers the `Release` workflow and publishes:

- GitHub release artifacts in `mrSamDev/brew-potato`
- `Casks/brew-potato.rb` in `mrSamDev/homebrew-tap`

### Install via Homebrew

```sh
brew tap mrSamDev/homebrew-tap
brew install --cask brew-potato
```

## Keybindings

| Key       | Action             |
| --------- | ------------------ |
| `↑` / `↓` | Navigate packages  |
| `d`       | Uninstall selected |
| `?`       | Show credits       |
| `q`       | Quit               |

## License

Do what every you want with it.
