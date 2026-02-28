# brew-ui

A polished terminal UI for managing Homebrew packages, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Full-screen TUI with rounded table border and styled header/footer
- Lists all user-installed Homebrew formulae with install date
- Async uninstall with a live spinner — row turns red in-place, no view flicker
- Graceful error display if `brew` is unavailable

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
go build -o brew-ui .
./brew-ui
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

## Keybindings

| Key       | Action             |
|-----------|--------------------|
| `↑` / `↓` | Navigate packages  |
| `d`       | Uninstall selected |
| `q`       | Quit               |
