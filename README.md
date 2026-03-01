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

## Next Steps

[GoReleaser](https://github.com/goreleaser/goreleaser) can publish brew-potato as an actual Homebrew formula so anyone can install it with `brew install`. It handles cross-compilation, GitHub release assets, and writing the formula to a tap repo automatically.

## Keybindings

| Key       | Action             |
| --------- | ------------------ |
| `↑` / `↓` | Navigate packages  |
| `d`       | Uninstall selected |
| `q`       | Quit               |

## License

Do what every you want with it.
